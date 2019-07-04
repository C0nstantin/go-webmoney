package wmsigner

/* example for use

1.s:=signer.Newsigner(wmid,pass,key)
2.return_value :=s.Sign(String)
*/

import (
	"fmt"
	"encoding/base64"
	"bytes"
	"encoding/binary"
	"golang.org/x/crypto/md4"
  "regexp"
	"errors"
	"math/big"
	"crypto/rand"
)
/*================= public package function ==========================*/

/*return wmsigner object for next using*/
func NewSigner(wmid string,pass string, key string )(*wmsigner){
	return &wmsigner{wmid,pass,key,nil}
}


/*================ private package function =========================*/

/*
function hashing  butes slice to md4
*/
func md4Hash(buff []byte)([]byte){
	hash:= md4.New()
	hash.Write(buff)
	return hash.Sum(nil);
}

/*
xor two byte slice
*/
func xor_wm(buff []byte,hash []byte ) ([]byte){
	b:=bytes.NewBuffer([]byte{});
	for i:=0; i<len(buff); i++ {
		b.WriteByte(buff[i]^hash[i%len(hash)])
		}
	return b.Bytes();
}

/*
reverse byte in bytes slice
*/

func reversebyte(s []byte) []byte{
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
	    s[i], s[j] = s[j], s[i]
	}
	return s;
}

/*create bignum object from bytes.Buffer*/
func readBn(b *bytes.Buffer )(*big.Int){
  var count uint16
	binary.Read(b,binary.LittleEndian,&count);
	data := reversebyte(b.Next(int(count)))
	i:=new(big.Int)
	i.SetBytes(data);
  return i
}

/*================= wmsigner object ==============================*/

type wmsigner struct{
	wmid string
	pass string
	key  string
	err  error
}

func (w *wmsigner) Sign(data string) (string,error){

	//check wmid
	if matched,_ := regexp.MatchString("[0-9]{12}",w.wmid); !matched{
     return "",errors.New("Wmid is incorrected")
	}

	//check key
	 decode_key, err := base64.StdEncoding.DecodeString(w.key);
	if err != nil {

		return "", err
	}
	if(len(decode_key)!=164){
	   return "",errors.New("Illegal size for keydata")
	}

	//checkpassword
	if(len(w.pass)==0){
   return "",errors.New("Password can not to be empty!")
	}
	powers,modules:=w.initkey()
	hash:=md4Hash([]byte(data))
	buf:=bytes.NewBuffer(hash)
  var ini uint32
	//add rendom 40 bytes
	for i:=1;i<=10;i++{
		binary.Read(rand.Reader, binary.LittleEndian, &ini)
		ini=0;
	  binary.Write(buf,binary.LittleEndian,ini)
	}
	hash=buf.Bytes()
	lenn:=make([]byte,2)
	binary.LittleEndian.PutUint16(lenn,uint16(len(hash)))
	b:=bytes.NewBuffer(lenn)
	b.Write(hash)

	for b.Len()<len(modules.Bytes()){
		b.WriteByte(byte(0));
	}


	data_res:=b.Bytes()
	data_res=reversebyte(data_res)
	i:=new(big.Int)
	i.SetBytes(data_res);
	signature:=i.Exp(i,powers,modules)

	sig:=bytes.NewBuffer([]byte{})
	for i:=len(signature.Bytes());i<=len(modules.Bytes());i++{
	sig.Write([]byte{0})
	}
	sig.Write(signature.Bytes())
	sigbyte:=reversebyte(sig.Bytes());
	sig=bytes.NewBuffer(sigbyte);
  var s = make([]uint16,len(sig.Bytes())/2)
	err=binary.Read(sig,binary.LittleEndian,s)
	if(err!=nil){
		return "", err;
	}
	result:=""
  for _,val:= range s {
	  result+=fmt.Sprintf("%04x",val);
	}

	return result ,nil

}

func (w *wmsigner) initkey() (*big.Int,*big.Int) {
  data,_  := base64.StdEncoding.DecodeString(w.key);

	r:=bytes.NewBuffer(data)
	header:=r.Next(30);
	body:=r.Bytes()
	body=w.encryptKey(body)
	r=bytes.NewBuffer([]byte{})

	r.Write(header)
	r.Write(body)

	reserved1:=r.Next(2)
	r.Next(2)
	crc:=r.Next(16)
	length:=r.Next(4)

	body=r.Bytes()
	r=bytes.NewBuffer([]byte{})
	reserved1=append(reserved1,byte(0),byte(0))
	r.Write(reserved1)
	r.Write(bytes.Repeat([]byte{0},16))
	r.Write(length)
	r.Write(body)

	calculated_crc:=md4Hash(r.Bytes())
	if !bytes.Equal(crc, calculated_crc) {
		if w.err == nil{
		w.pass=w.pass[0:len(w.pass)/2]
		w.err=errors.New("invalid crc")
	  return  w.initkey();
		}else{
	    fmt.Println(w.err);
		  return nil, nil

		}
	}
	r=bytes.NewBuffer(body)
	r.Next(4);
	powers:=readBn(r);
	modules:=readBn(r);
	return powers,modules

}

func(w *wmsigner) encryptKey(buff []byte) ([]byte){
	hashResult:=md4Hash([]byte(w.wmid+w.pass))
	return xor_wm(buff,hashResult)
}

func add_tobegin(b []byte, bb []byte) {
	buf:=bytes.NewBuffer(bb);
	buf.Write(b)
	b=buf.Bytes()
}

