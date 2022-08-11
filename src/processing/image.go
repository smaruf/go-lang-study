package main

import (
   "log"
   "os"
)

type MosaicBuilder struct {
   logger *log.Logger
}

func NewMosaicBuilder() *MosaicBuilder {
   return &MosaicBuilder{
      logger: log.New(os.Stderr, "[mosaic] ", log.LstdFlags),
   }
}

func (b *MosaicBuilder) Build() error {
   return nil
}

func main() {
   builder := NewMosaicBuilder()

   if err := builder.Build(); err != nil {
      builder.logger.Fatal(err)
   }
}
