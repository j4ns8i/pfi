use std::error::Error;

use clap::Parser;
use image::io::Reader;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
struct Pfi {
    image: String,
}

fn main() {
    let pfi = Pfi::parse();
    if let Err(e) = pfi.run() {
        println!("Error: {}", e);
        std::process::exit(1)
    }
}

impl Pfi {
    fn run(&self) -> Result<(), Box<dyn Error>> {
        let img = Reader::open(&self.image)?.decode()?;
        img.save("asdf.png")?;

        Ok(())
    }
}
