data "libp2p_peer_id" "this" {
  ed25519_secret_key = base64encode("00000000000000000000000000000001")
}
