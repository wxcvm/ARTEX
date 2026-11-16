package report

import (
	"bytes"
	"errors"
	"image"
)

const MaxReportLogoBytes = 2 << 20

func reportLogoConfig(raw []byte)(image.Config,string,error){
	if len(raw)==0||len(raw)>MaxReportLogoBytes{return image.Config{},"",errors.New("logo must be a PNG or JPEG up to 2 MiB")}
	cfg,format,err:=image.DecodeConfig(bytes.NewReader(raw))
	if err!=nil||(format!="png"&&format!="jpeg")||cfg.Width<=0||cfg.Height<=0||cfg.Width>4096||cfg.Height>4096{
		return image.Config{},"",errors.New("logo must be a valid PNG or JPEG up to 4096 pixels per side")
	}
	return cfg,format,nil
}

func ValidateReportLogo(raw []byte)error{
	if _,_,err:=reportLogoConfig(raw);err!=nil{return err}
	if _,_,err:=image.Decode(bytes.NewReader(raw));err!=nil{return errors.New("invalid or truncated logo image")}
	return nil
}
