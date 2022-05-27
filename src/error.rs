use std::fmt;
use std::error::Error;

#[derive(Debug)]
pub enum PfiError {
    Io(std::io::Error),
    ImageError(image::ImageError),
}

impl std::error::Error for PfiError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match self {
            Self::Io(e) => Some(e),
            Self::ImageError(e) => Some(e),
        }
    }
}

impl fmt::Display for PfiError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Self::Io(ref e) => write!(f, "{}", e),
            Self::ImageError(ref e) => write!(f, "{}", e),
        }
    }
}

impl From<std::io::Error> for PfiError {
    fn from(error: std::io::Error) -> Self {
        PfiError::Io(error)
    }
}

impl From<image::ImageError> for PfiError {
    fn from(error: image::ImageError) -> Self {
        PfiError::ImageError(error)
    }
}
