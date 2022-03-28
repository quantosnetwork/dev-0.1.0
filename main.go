/*
Copyright © 2022 Quantos Developers <dev@quantos.network>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"github.com/quantosnetwork/dev-0.1.0/cmd"
	_ "github.com/quantosnetwork/dev-0.1.0/config"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
	/*buf := make([]byte, 32)
	frand.Read(buf)
	account.CreateNewAddress(byte(0x00), hex.EncodeToString(buf)[:32])*/
	select {}

}

// methods override
