use kmeans::{KMeans, KMeansConfig};
use termion::color;

pub type Color24Bit = image::Rgb<u8>;

/// Constructs a palette of num_colors from the given img.
pub fn palette_from_img(img: image::RgbImage, num_colors: usize) -> Vec<Color24Bit> {
    let km = KMeans::new(
        normalize_to_f32(img.to_vec()), // fed in as [r, g, b, r, g, b, ...]
        (img.width() * img.height()) as usize,
        3, // dimensions
    );

    let clusters = km.kmeans_lloyd(
        num_colors,
        100, // iterations
        KMeans::init_kmeanplusplus,
        &KMeansConfig::default(),
    );

    let normalized_u8 = normalize_to_u8(clusters.centroids);
    let rgb_chunks = normalized_u8.chunks_exact(3);
    rgb_chunks.map(|x| image::Rgb([x[0], x[1], x[2]])).collect()
}

/// Converts integers within [0, 255] to floating point values within [0.0, 1.0].
fn normalize_to_f32(uints: Vec<u8>) -> Vec<f32> {
    uints
        .into_iter()
        .map(|v| (v as f32 / u8::MAX as f32))
        .collect()
}

/// Converts floating point values expected to be within [0.0, 1.0] into integers in [0, 255].
fn normalize_to_u8(floats: Vec<f32>) -> Vec<u8> {
    floats
        .into_iter()
        .map(|x| (f32::from(u8::MAX) * x) as u8)
        .collect()
}

/// Creates a block (two spaces) with the given background color.
pub fn color_block(color: &Color24Bit) -> String {
    format!(
        "{}  {}",
        color::Bg(color::Rgb(color[0], color[1], color[2])),
        color::Bg(color::Reset)
    )
}

#[cfg(test)]
mod tests {
    #[test]
    fn normalize_to_f32() {
        let input = vec![255, 0, 255, 0];
        let expected = [1.0, 0.0, 1.0, 0.0];
        let normalized = super::normalize_to_f32(input);
        let all_equal = normalized
            .into_iter()
            .zip(expected)
            .all(|(e, n)| (e - n).abs() <= f32::EPSILON);
        assert!(all_equal);
    }

    #[test]
    fn normalize_to_u8() {
        let input = vec![1.0, 0.5, 0.0, 1.0];
        let expected = vec![255, 127, 0, 255];
        let normalized = super::normalize_to_u8(input);
        assert_eq!(normalized, expected);
    }
}
