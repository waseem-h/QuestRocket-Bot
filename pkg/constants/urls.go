package constants

const NYCPokemap = "https://nycpokemap.com/"
const VANPokemap = "https://vanpokemap.com/"
const SGPokemap = "https://sgpokemap.com/"
const SYDNEYPogomap = "https://sydneypogomap.com/"

func URLMap() map[string]string {
    return map[string]string{
        "NewYork Pokemap":   NYCPokemap,
        "Vancouver Pokemap": VANPokemap,
        "Singapore Pokemap": SGPokemap,
        "Sydney Pogomap":    SYDNEYPogomap,
    }
}
