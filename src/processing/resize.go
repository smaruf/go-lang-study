func (b *MosaicBuilder) loadImage(path string) (image.Image, error) {
   infile, err := os.Open(path)
   defer func(infile *os.File) {
      _ = infile.Close()
   }(infile)

   if err != nil {
      return nil, err
   }

   src, err := tga.Decode(infile)
   if err != nil {
      return nil, err
   }
   
   return resize.Resize(b.partSize, b.partSize, src, resize.Lanczos3), nil
}
