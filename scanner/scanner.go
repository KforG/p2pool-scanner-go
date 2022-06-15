package scanner

import (
	"fmt"
	"time"

	"github.com/KforG/p2pool-scanner-go/config"
	"github.com/KforG/p2pool-scanner-go/geo"
	"github.com/KforG/p2pool-scanner-go/logging"
	"github.com/KforG/p2pool-scanner-go/node"
	"github.com/KforG/p2pool-scanner-go/util"
)

type Node struct {
	IP          string           `json:"ip"`
	DomainName  string           `json:"domain_name"`
	LocalStats  node.NodeStats   `json:"node_stats"`
	GlobalStats node.GlobalStats `json:"global_stats"`
	GeoLocation geo.Geo          `json:"geo_location"`
}

type Nodes []Node

func Scanner(n *Nodes) {
	go node.GetAllDomainIPs()
	n.loadBootstrapNodes()

	for {
		logging.Infof("Scanning P2pool network for public P2Pool nodes...\n")
		for i := 0; i < len(*n); i++ {
			err := node.GetLocalStats((*n)[i].IP, &(*n)[i].LocalStats)
			if err != nil {
				logging.Infof("%s, did not return updated node state, might be unreachable. Removing from node list.\n", (*n)[i].IP)
				// bring element to remove at the end if its not there yet
				if i != len(*n)-1 {
					(*n)[i] = (*n)[len(*n)-1]
				}
				// drop the last element
				*n = (*n)[:len(*n)-1]

				continue
			}

			peers, _ := node.GetPeers((*n)[i].IP)
			go n.discoverNewNodes(peers)

			_ = node.GetGlobalStats((*n)[i].IP, &(*n)[i].GlobalStats)
		}
		if config.Active.RescanTime > 0 {
			logging.Infof("Updating nodes stats again and rescanning for new peers in %d minutes", config.Active.RescanTime)
			time.Sleep(time.Duration(config.Active.RescanTime) * time.Minute)
		}
	}
}

func (n *Nodes) loadBootstrapNodes() {
	logging.Infof("Loading bootstrap nodes..")
	for i := 0; i < len(config.Active.BootstrapNodes); i++ {
		var bn Node
		err := node.GetLocalStats(config.Active.BootstrapNodes[i], &bn.LocalStats)
		if err != nil {
			logging.Infof("%s couldn't be reached, consider updating bootstrap nodes\n", config.Active.BootstrapNodes[i])
			continue
		}
		err = node.GetGlobalStats(config.Active.BootstrapNodes[i], &bn.GlobalStats)
		if err != nil {
			logging.Errorf("Unable to fetch GlobalStats from %s\n", config.Active.BootstrapNodes[i])
		}

		bn.DomainName = config.Active.BootstrapNodes[i]
		node.GetIPFromDomain(bn.DomainName, &bn.IP)

		err = geo.GetGeoLocation(bn.IP, &bn.GeoLocation)
		if err != nil {
			logging.Errorf("Unable to fetch geolocation of %s:%s\n", bn.DomainName, bn.IP)
		}
		*n = append(*n, bn)
	}
}

func (n *Nodes) discoverNewNodes(peers []string) {
	n.removeDuplicatePeers(&peers) // remove possible peer duplicates
	for i := 0; i < len(peers); i++ {
		var an Node

		err := node.GetLocalStats(fmt.Sprintf(peers[i]+":"+config.Active.Port), &an.LocalStats)
		if err != nil {
			// If the peer doesn't respond to this request, it is assumed inactive/private
			continue
		}

		err = node.GetGlobalStats(fmt.Sprintf(peers[i]+":"+config.Active.Port), &an.GlobalStats)
		if err != nil {
			logging.Errorf("Failed to fetch Global stats from %s, node seems to be alive\n", fmt.Sprintf(peers[i]+":"+config.Active.Port), err)
		}

		err = geo.GetGeoLocation(peers[i], &an.GeoLocation)
		if err != nil {
			logging.Errorf("Failed to get Geolocation of %s\n", peers[i])
		}

		an.IP = fmt.Sprintf(peers[i] + ":" + config.Active.Port)
		node.CheckForDomain(an.IP, &an.DomainName)

		if n.checkForDuplicatePeer(an.IP) {
			*n = append(*n, an)
		}
	}
}

func (n *Nodes) removeDuplicatePeers(peers *[]string) {
	for i := 0; i < len(*n); i++ {
		for j := 0; j < len(*peers); {
			if (*n)[i].IP == fmt.Sprintf((*peers)[j]+":"+config.Active.Port) {
				// peer is the same as already known node
				// needs to be removed from slice of peers
				util.RemoveStringSliceIndex(j, peers)
				continue
			}
			j++
		}
	}
}

// This function is used to make sure that a new node hasn't been added by another go routine.
// If it's not here possible duplicates will appear.
// returns true if peer is new
func (n *Nodes) checkForDuplicatePeer(newNode string) bool {
	for i := 0; i < len(*n); i++ {
		if (*n)[i].IP == newNode {
			return false
		}
	}
	return true
}
