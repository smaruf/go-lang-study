func (b *MosaicBuilder) Build() error {
   b.logger.Println("Load parts paths ...")
   partsPaths, err := b.getPartsPaths()
   if err != nil {
      return err
   }

   b.logger.Println("Load parts map ...")
   partsMap, err := b.getPartsMap(partsPaths)
   if err != nil {
      return err
   }
   
   return nil
}
