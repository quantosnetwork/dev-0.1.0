package tests

import (
	"encoding/hex"
	"github.com/quantosnetwork/dev-0.1.0/color"
	"github.com/quantosnetwork/dev-0.1.0/core/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type AccountAddressTestSuite struct {
	suite.Suite

	Name    string
	Key     string
	Address string
	NetID   byte
	Path    string
}

func (suite *AccountAddressTestSuite) SetupTest() {
	suite.Key = "qwerty1234"
	suite.NetID = 0x00
	suite.Path = "./data/.keys/"
	suite.Name = "AccountAddressTestSuite"
}

func (suite *AccountAddressTestSuite) TestAccountAddress() {
	var err error
	suite.Address, err = account.CreateNewAddress(suite.NetID, suite.Key)
	if assert.NoError(suite.T(), err, "address generated ok!") { //nolint:typecheck
		log.Printf(color.Blue+"address was generated! current address: %s \n", suite.Address)
		if assert.NotEmpty(suite.T(), suite.Address, "address is not empty!") { //nolint:typecheck
			log.Println(color.Green + "pass: address is not empty")
		}
		if !assert.FileExists(suite.T(), suite.Path+suite.Address+".wal", "lock file for address exists!") { //nolint:typecheck
			log.Println(color.Red + "address lock file does not exists!")
		} else {
			log.Println(color.Green + "address was successfully encrypted and locked on disk")
		}
	}
}

func (suite *AccountAddressTestSuite) TestAccountAddressDecrypt() {

	decrypt, err := account.GetAddressFromStorage(suite.Address, suite.Key)
	if assert.NoError(suite.T(), err, "an error occured: %v") { //nolint:typecheck

		if assert.Equal(suite.T(), hex.EncodeToString(decrypt.PubKey), suite.Address) { //nolint:typecheck
			log.Println(color.Blue + "decryption test succeeded addresses are the same!")

		} else {
			suite.Fail(color.Red + "both addresses are not equal! wrong decryption") //nolint:typecheck
		}
	} else {
		suite.Fail(color.Red + "an error occured while decrypting, wrong key?") //nolint:typecheck
	}

}

func (suite *AccountAddressTestSuite) TestAccountAddressToString() {

	addrType, _ := hex.DecodeString(suite.Address)
	addrWithPrefix := account.Address(addrType)
	if assert.Contains(suite.T(), addrWithPrefix.String(), "QBX") { //nolint:typecheck
		log.Printf(color.Blue+"address is having the right prefix (QBX): %s", addrWithPrefix.String())
	} else {
		suite.Fail(color.Red + "address is invalid (does not have the prefix QBX)") //nolint:typecheck
	}
}

func (suite *AccountAddressTestSuite) BeforeTest(suitename, testname string) {
	log.Println(suitename + "--" + testname)
	log.Printf("--------------------------------------------------------------\n\n")
}

func (suite *AccountAddressTestSuite) AfterTest(suitename, testname string) {
	log.Printf("--------------------------------------------------------------\n\n")

	log.Println(color.Green + "(" + testname + ") PASS")

}

func TestAccountAddressTestSuite(t *testing.T) {
	suite.Run(t, new(AccountAddressTestSuite))
}
