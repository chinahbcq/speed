package service

import (
    "utils"
)

type producer struct {}
var Producer producer

func (handle *producer) Push(input []byte, logbuf *utils.LogBuffer) string {
    //do something
    
    return string(input)
} 
