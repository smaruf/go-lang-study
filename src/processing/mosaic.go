type MosaicBuilder struct {
   pathToParts string
   partSize    uint
   logger      *log.Logger
}

func (b *MosaicBuilder) getPartsPaths() ([]string, error) {
   var res []string

   return res, filepath.Walk(b.pathToParts, func(path string, info os.FileInfo, err error) error {
      if err != nil {
         return err
      }

      if !info.IsDir() {
         res = append(res, path)
      }

      return nil
   })
}
