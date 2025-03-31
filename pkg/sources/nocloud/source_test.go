package nocloud_test

var networkConfigResponse = `version: 2
ethernets:
  eno1:
    renderer: networkd
    nameservers:
      addresses:
        - 1.1.1.1
        - 8.8.8.8
    match:
      macaddress: "26:3c:0b:14:12:11"
    dhcp4: false
    addresses:
      - 100.100.100.2/24
    routes:
      - to: 0.0.0.0/0
        via: 100.100.100.1

vlans:
  vlan6:
    id: 6
    link: eno1
    dhcp4: false
    addresses:
      - 90.90.90.2/16
    nameservers:
      addresses:
        - 1.1.1.1
        - 8.8.8.8
    routes:
      - to: 0.0.0.0/0
        via: 90.90.90.1
        metric: 1000
`

var metadataResponse = `
instance-id: host0001
local-hostname: host0001
`

var userDataResponse = `#cloud-config
hostname: host0001

disable_root: false
ssh_pwauth: false

users:
  - name: root
    ssh_authorized_keys:
        - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDspDiCNcycfJwGU8AnX+6+SU9l/EMVclMreFJL2RLa6XU56lVWIk7s/xl32C/UpTo1BT9CRmWYeJQtqGbErGtmTLHzqnV03WanZgHViOi8LoQvt4bd1ssCGIbohPHbmaGaiZiDqtYCQb4NiweVyIK14qUQDawdLQ3CLcMPVzcByIb/mvb/fBouQkqeNZHURBp+40o1CGCqsu3gIMhP6Pd/ncF9p71eCvZ+cRkBuAqpOVnxN7f4dW/Bt66imHX8BZCoQBMenpe6CNhVTkmQQbVmHBenK5Err2IYkVsgHzcLvhxAr/tsKh2rkly9HZYGF+6xemKzAAH+xp3wUrYWlLHIYMXyCOXL1EHqSjcub/i/cjxbgssUqJuwBTrjV5bl9VP2y+4MlE3YM+JvPZebl604GXytT0RAk35jqeU+4gh616CaykdvXyLNPBL4tsHqUy7IQzvvC83SdtYroMmqnDp4uXKYwi8UPEa17jBOIsKtRZLWa4vr0wCUbnMV0Nj3kCE=

chpasswd:
  expire: false

runcmd:
  - userdel -r ubuntu || true
  - systemctl restart ssh
`
