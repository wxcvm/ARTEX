package server

import (
	"net/http"
	"sort"
)

type purgeIntelAnalysis struct {
	Schema string `json:"schema"`
	Target string `json:"target"`
	HistoricalOnly bool `json:"historical_only"`
	Truncated bool `json:"truncated"`
	ExcludedDNS int `json:"excluded_dns"`
	Names int `json:"names"`
	IPs int `json:"ips"`
	Edges int `json:"edges"`
	PortObservations int `json:"port_observations"`
	SourceCounts map[string]int `json:"source_counts"`
	NamesWithMultipleSources []string `json:"names_with_multiple_sources"`
	NamesWithoutDNS []string `json:"names_without_dns"`
	NamesWithDNSFailures []string `json:"names_with_dns_failures"`
	NamesWithoutRetainedAddress []string `json:"names_without_retained_address"`
	IPsWithoutDomainLink []string `json:"ips_without_domain_link"`
	IPsWithoutPortEvidence []string `json:"ips_without_port_evidence"`
	UnknownTransportPorts int `json:"unknown_transport_ports"`
	PortsWithoutObservationTime int `json:"ports_without_observation_time"`
	SharedAddresses int `json:"shared_addresses"`
}

func analyzePurgeIntel(b purgeIntelBundle) purgeIntelAnalysis {
	a := purgeIntelAnalysis{
		Schema: "artex.purge-intel-analysis.v1", Target:b.Target, HistoricalOnly:true,
		Truncated:b.Truncated, ExcludedDNS:b.Excluded["dns_outside_scope"],
		Names:len(b.Names), IPs:len(b.IPs), Edges:len(b.Edges),
		SourceCounts:map[string]int{},
		NamesWithMultipleSources:[]string{}, NamesWithoutDNS:[]string{},
		NamesWithDNSFailures:[]string{}, NamesWithoutRetainedAddress:[]string{},
		IPsWithoutDomainLink:[]string{}, IPsWithoutPortEvidence:[]string{},
	}
	linked := map[string]map[string]bool{}
	for _, edge := range b.Edges {
		if edge.Relation != "A" && edge.Relation != "AAAA" { continue }
		if linked[edge.To] == nil { linked[edge.To] = map[string]bool{} }
		linked[edge.To][edge.From] = true
	}
	for _, n := range b.Names {
		sources := map[string]bool{}
		for _, e := range n.Evidence {
			if e.Source != "registered-scope" { sources[e.Source] = true }
		}
		for source := range sources { a.SourceCounts[source]++ }
		if len(sources)>1 { a.NamesWithMultipleSources=append(a.NamesWithMultipleSources,n.Name) }
		if n.DNSCollectedAt==nil { a.NamesWithoutDNS=append(a.NamesWithoutDNS,n.Name) }
		if len(n.DNSErrors)>0 {
			for _, e := range n.DNSErrors {
				// Absent record types are normal, not resolver failures.
				if e.Code!="ENODATA" && e.Code!="ENOTFOUND" {
					a.NamesWithDNSFailures=append(a.NamesWithDNSFailures,n.Name);break
				}
			}
		}
		hasAddr:=false
		for _, d:=range n.DNS { if d.Type=="A"||d.Type=="AAAA" {hasAddr=true;break} }
		if !hasAddr { a.NamesWithoutRetainedAddress=append(a.NamesWithoutRetainedAddress,n.Name) }
	}
	for _, ip:=range b.IPs {
		if len(linked[ip.IP])==0 { a.IPsWithoutDomainLink=append(a.IPsWithoutDomainLink,ip.IP) }
		if len(linked[ip.IP])>1 { a.SharedAddresses++ }
		if len(ip.PassivePorts)==0 { a.IPsWithoutPortEvidence=append(a.IPsWithoutPortEvidence,ip.IP) }
		for _, p:=range ip.PassivePorts {
			a.PortObservations++
			if p.Proto=="unknown" {a.UnknownTransportPorts++}
			timed:=false
			sources:=map[string]bool{}
			for _, e:=range p.Evidence {
				sources[e.Source]=true
				if e.ObservedAt!=nil && *e.ObservedAt!="" {timed=true}
			}
			for source:=range sources {a.SourceCounts[source]++}
			if !timed {a.PortsWithoutObservationTime++}
		}
	}
	for _, list:=range [][]string{a.NamesWithMultipleSources,a.NamesWithoutDNS,a.NamesWithDNSFailures,a.NamesWithoutRetainedAddress,a.IPsWithoutDomainLink,a.IPsWithoutPortEvidence} {sort.Strings(list)}
	return a
}

func (s *Server) getPurgeIntelAnalysis(w http.ResponseWriter,r *http.Request) {
	b,ok:=s.loadPurgeIntelSnapshot(w,r)
	if !ok {return}
	writeJSON(w,http.StatusOK,analyzePurgeIntel(b))
}
