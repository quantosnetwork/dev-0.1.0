package main

import (
	"encoding/hex"
	"github.com/quantosnetwork/v0.1.0-dev/account"
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
		log.Printf("address was generated! current address: %s \n", suite.Address)
		if assert.NotEmpty(suite.T(), suite.Address, "address is not empty!") { //nolint:typecheck
			log.Println("pass: address is not empty")
		}
		if !assert.FileExists(suite.T(), suite.Path+suite.Address+".wal", "lock file for address exists!") { //nolint:typecheck
			log.Println("address lock file does not exists!")
		} else {
			log.Println("address was successfully encrypted and locked on disk")
		}
	}
}

func (suite *AccountAddressTestSuite) TestAccountAddressDecrypt() {

	decrypt, err := account.GetAddressFromStorage(suite.Address, suite.Key)
	if assert.NoError(suite.T(), err, "an error occured: %v") { //nolint:typecheck

		if assert.Equal(suite.T(), hex.EncodeToString(decrypt.PubKey), suite.Address) { //nolint:typecheck
			log.Println("decryption test succeeded addresses are the same!")

		} else {
			suite.Fail("both addresses are not equal! wrong decryption") //nolint:typecheck
		}
	} else {
		suite.Fail("an error occured while decrypting, wrong key?") //nolint:typecheck
	}

}

func (suite *AccountAddressTestSuite) TestAccountAddressToString() {

	addrType, _ := hex.DecodeString(suite.Address)
	addrWithPrefix := account.Address(addrType)
	if assert.Contains(suite.T(), addrWithPrefix.String(), "QBX") { //nolint:typecheck
		log.Printf("(ADDRESS FORMAT AND PREFIX):  address is having the right prefix (QBX): %s", addrWithPrefix.String())
	} else {
		suite.Fail("(ADDRESS FORMAT AND PREFIX): address is invalid (does not have the prefix QBX)") //nolint:typecheck
	}
}

func (suite *AccountAddressTestSuite) BeforeTest(suitename, testname string) {
	log.Println(suitename + "--" + testname)
	log.Printf("--------------------------------------------------------------\n\n")
}

func (suite *AccountAddressTestSuite) AfterTest(suitename, testname string) {
	log.Printf("--------------------------------------------------------------\n\n")

	log.Println("RESULT: PASS")

}

func TestAccountAddressTestSuite(t *testing.T) {
	suite.Run(t, new(AccountAddressTestSuite))
}
