package main

import "openspy.net/natneg-helper/src/Messages"

type INatNegCore interface {
	HandleInitMessage(msg Messages.InitMessage)
}
