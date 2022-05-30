package zdpgo_zip

import (
	"github.com/zhangdapeng520/zdpgo_test/assert"
	"testing"
)

/*
@Time : 2022/5/30 14:22
@Author : 张大鹏
@File : zip_test.go
@Software: Goland2021.3.1
@Description:
*/

func getZip() *Zip {
	return NewWithConfig(&Config{
		Debug: true,
	})
}

func TestZip_Zip(t *testing.T) {
	z := getZip()
	err := z.Zip("test2", "test.zip")
	assert.NoError(t, err)
}

func TestZip_ZipAndDelete(t *testing.T) {
	z := getZip()
	err := z.ZipAndDelete("test3", "test3.zip")
	assert.NoError(t, err)
}

func TestZip_Unzip(t *testing.T) {
	z := getZip()
	err := z.Unzip("test3.zip", "test3")
	assert.NoError(t, err)
}

func TestZip_UnzipToCurrentDir(t *testing.T) {
	z := getZip()
	err := z.UnzipToCurrentDir("test.zip")
	assert.NoError(t, err)
}

func TestZip_UnzipToCurrentDirAndDelete(t *testing.T) {
	z := getZip()
	err := z.UnzipToCurrentDirAndDelete("test.zip")
	assert.NoError(t, err)
}
