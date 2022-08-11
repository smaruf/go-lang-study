func (b *MosaicBuilder) getPartsMap(parts []string) (map[PixelColor]image.Image, error) {
   partsMap := make(map[PixelColor]image.Image, len(parts))

   for _, path := range parts {
      src, err := b.loadImage(path)
      if err == nil {
         partsMap[calculateModalAverageColour(src)] = src
      }
   }

   if len(partsMap) == 0 {
      return nil, fmt.Errorf("empty map")
   }

   return partsMap, nil
}
