package node

import (
	"fmt"
	"net"
	"strings"

	"github.com/KforG/p2pool-scanner-go/config"
	"github.com/KforG/p2pool-scanner-go/logging"
	"github.com/KforG/p2pool-scanner-go/util"
)

func GetPeers(IP string) (peers []string, err error) {
	var jsonPayload string
	err = util.GetJson(fmt.Sprintf("%s/peer_addresses", IP), &jsonPayload)
	if err != nil {
		return peers, err
	}
	//Since the http response is a string we iterate over the string to split up connected peers and remove potential port number
	//SplitN from the string package can seperate the peers, only need to iterate over the slice to remove potential port number
	peerSplit := strings.SplitN(jsonPayload, " ", len(jsonPayload))
	for i := 0; i < len(peerSplit); i++ {
		k := strings.SplitN(peerSplit[i], ":", len(peerSplit))
		peers = append(peers, k[0])
	}
	return peers, err
}

type NodeStats struct {
	EfficiencyIfMinerPerfect float64 `json:"efficiency_if_miner_perfect"`
	Efficiency               float64 `json:"efficiency"`
	Fee                      float64 `json:"fee"`
	DonationProportion       float64 `json:"donation_proportion"`
	Uptime                   float64 `json:"uptime"`
	BlockValue               float64 `json:"block_value"`
	AttemptsToBlock          uint64  `json:"attempts_to_block"`
	//MinerHashRates               map[string]string            `json:"miner_hash_rates"`		// Wrong type. Doesn't matter, isn't useful.
	//MinerDeadHashRates           map[string]string            `json:"miner_dead_hash_rates"`
	//MinerLastDifficulties        map[string]string            `json:"miner_last_difficulties"`
	Peers                        Peers                        `json:"peers"`
	MyShareCountsInLastHour      MyShareCountsInLastHour      `json:"my_share_counts_in_last_hour"`
	MyStaleProportionsInLastHour MyStaleProportionsInLastHour `json:"my_stale_proportions_in_last_hour"`
	Shares                       Shares                       `json:"shares"`
	MyHashRatesInLastHour        MyHashRatesInLastHour        `json:"my_hash_rates_in_last_hour"`
	Version                      string                       `json:"version"`
	Warnings                     []string                     `json:"warnings"`
	ProtocolVersion              uint                         `json:"protocol_version"`
	AttemptsToShare              uint                         `json:"attempts_to_share"`
}
type Peers struct {
	Outgoing uint `json:"outgoing"`
	Incoming uint `json:"incoming"`
}
type MyShareCountsInLastHour struct {
	DoaStaleShares    uint `json:"doa_stale_shares"`
	StaleShares       uint `json:"stale_shares"`
	OrphanStaleShares uint `json:"orphan_stale_shares"`
	Shares            uint `json:"shares"`
	UnstaleShares     uint `json:"unstale_shares"`
}
type MyStaleProportionsInLastHour struct {
	OrphanStale float64 `json:"orphan_stale"`
	Stale       float64 `json:"stale"`
	DeadStale   float64 `json:"dead_stale"`
}
type Shares struct {
	Total  uint `json:"total"`
	Orphan uint `json:"orphan"`
	Dead   uint `json:"dead"`
}
type MyHashRatesInLastHour struct {
	Note     string  `json:"note"`
	Actual   float64 `json:"actual"`
	Rewarded float64 `json:"rewarded"`
	Nonstale float64 `json:"nonstale"`
}

func GetLocalStats(IP string, node *NodeStats) error {
	err := util.GetJson(fmt.Sprintf("%s/local_stats", IP), node)
	if err != nil {
		return err
	}
	return nil
}

type GlobalStats struct {
	PoolNonstaleHashRate   float64 `json:"pool_nonstale_hash_rate"`
	PoolStaleProp          float64 `json:"pool_stale_prop"`
	PoolHashRate           float64 `json:"pool_hash_rate"`
	NetworkHashrate        float64 `json:"network_hashrate"`
	NetworkBlockDifficulty float64 `json:"network_block_difficulty"`
	MinDifficulty          float64 `json:"min_difficulty"`
}

func GetGlobalStats(url string, globStats *GlobalStats) error {
	err := util.GetJson(fmt.Sprintf("%s/global_stats", url), globStats)
	if err != nil {
		return err
	}
	return nil
}

// This function does a DNS lookup and assigns the first known IP to the domain name
func GetIPFromDomain(domainWithPort string, IP *string) {
	domain := strings.SplitN(domainWithPort, ":", len(domainWithPort))
	IPs, err := net.LookupIP(domain[0])
	if err != nil {
		logging.Errorf("Failed to get IP from %s\n", domain[0], err)
		*IP = domainWithPort
		return
	}
	*IP = fmt.Sprintf("%s:%s", IPs[0], domain[1]) // If bootstrapnode has other port than default we can't use the config variable
}

func GetAllDomainIPs() {
	for i := 0; i < len(config.Active.KnownDomains.NodeDomain); i++ {
		go func(i int) {
			GetIPFromDomain(config.Active.KnownDomains.NodeDomain[i].DomainName, &config.Active.KnownDomains.NodeDomain[i].IP)
		}(i)
	}
}

// The purpose of this func is to check if the IP of a node is associated
// with a known domain name.
// Current methods of checking databases for domains associated with IP addresses
// often return weird ISP strings.
// So we do some manual lookup in config
func CheckForDomain(IP string, domain *string) {
	if !config.Active.KnownDomains.Check {
		return
	}
	for i := 0; i < len(config.Active.KnownDomains.NodeDomain); i++ {
		if IP == config.Active.KnownDomains.NodeDomain[i].IP {
			*domain = config.Active.KnownDomains.NodeDomain[i].DomainName
			return
		}
	}
	*domain = IP
}
