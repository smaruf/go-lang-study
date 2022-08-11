type PixelColor [3]uint8

func calculateModalAverageColour(img image.Image) PixelColor {
   imgSize := img.Bounds().Size()

   var redTotal, greenTotal, blueTotal, pixelsCount int64

   for x := 0; x < imgSize.X; x++ {
      for y := 0; y < imgSize.Y; y++ {
         cc := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

         redTotal += int64(cc.R)
         greenTotal += int64(cc.G)
         blueTotal += int64(cc.B)

         pixelsCount++
      }
   }

   r := uint8(redTotal / pixelsCount)
   g := uint8(greenTotal / pixelsCount)
   b := uint8(blueTotal / pixelsCount)

   return PixelColor{r, g, b}
}
