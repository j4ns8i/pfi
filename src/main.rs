use clap::Parser;

mod error;
mod palette;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
struct Pfi {
    /// Image used to generate palette
    image: String,

    /// Number of colors in palette
    #[clap(short = 'n', default_value = "8")]
    num_clusters: usize,
}

fn main() {
    let p = Pfi::parse();
    if let Err(e) = run(p) {
        eprintln!("{}", e);
        std::process::exit(1);
    }
}

fn run(prog: Pfi) -> Result<(), error::PfiError> {
    let img = image::io::Reader::open(&prog.image)?.decode()?;
    let thumbnail = img.resize(500, 500, image::imageops::Nearest);
    let colors = palette::palette_from_img(thumbnail.into_rgb8(), prog.num_clusters);
    let line_of_colors = colors
        .iter()
        .map(palette::color_block)
        .collect::<Vec<String>>()
        .join(" ");
    println!("{}", line_of_colors);
    Ok(())
}
