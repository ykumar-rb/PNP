package templates

//template for preseed.cfg
var PreseedTmlp = `
choose-mirror-bin mirror/http/proxy string
d-i base-installer/kernel/override-image string linux-server
d-i netcfg/hostname string ztp-node-client
d-i netcfg/get_hostname string ztp-node-client
d-i clock-setup/utc boolean true
d-i clock-setup/utc-auto boolean true
d-i finish-install/reboot_in_progress note
d-i grub-installer/only_debian boolean true
d-i grub-installer/with_other_os boolean true
d-i     mirror/country          string United States
d-i mirror/http/mirror select CC.archive.ubuntu.com
#d-i mirror/http/mirror select 192.168.63.230/ubuntu/
d-i partman-efi/non_efi_system boolean false
d-i partman-auto-lvm/guided_size string max
d-i partman-auto/choose_recipe select atomic
d-i partman-auto/method string lvm
d-i partman-lvm/confirm boolean true
d-i partman-lvm/confirm boolean true
d-i partman-lvm/confirm_nooverwrite boolean true
d-i partman-lvm/device_remove_lvm boolean true
d-i partman/choose_partition select finish
d-i partman/confirm boolean true
d-i partman/confirm_nooverwrite boolean true
d-i partman/confirm_write_new_label boolean true
d-i pkgsel/include string openssh-server cryptsetup build-essential libssl-dev libreadline-dev zlib1g-dev linux-source dkms nfs-common
d-i pkgsel/install-language-support boolean false
d-i pkgsel/update-policy select none
d-i pkgsel/upgrade select full-upgrade
d-i time/zone string UTC
tasksel tasksel/first multiselect standard, ubuntu-server
#tasksel tasksel/first multiselect ubuntu-desktop

d-i console-setup/ask_detect boolean false
d-i keyboard-configuration/layoutcode string us
d-i keyboard-configuration/modelcode string pc105
d-i debian-installer/locale string en_US

# Create vagrant user account.
d-i passwd/user-fullname string vagrant
d-i passwd/username string vagrant
d-i passwd/user-password password vagrant
d-i passwd/user-password-again password vagrant
d-i user-setup/allow-password-weak boolean true
d-i user-setup/encrypt-home boolean false
d-i passwd/user-default-groups vagrant sudo
d-i passwd/user-uid string 900

d-i preseed/late_command string \
 in-target mkdir -P /root/PnP/certs/ ; \
 in-target wget -P /root/PnP/certs/ http://{{.IP}}:{{.MatchboxPort}}/assets/coreos/client/server.crt ; \
 in-target wget -P /root/PnP/ http://{{.IP}}:{{.MatchboxPort}}/assets/coreos/client/client ; \
 in-target wget -P /root/PnP/ http://{{.IP}}:{{.MatchboxPort}}/assets/coreos/client/bootstrap.sh ; \
 in-target chmod +x /root/PnP/bootstrap.sh; \
 in-target wget -P /root/PnP/ http://{{.IP}}:{{.MatchboxPort}}/assets/coreos/client/resolv.conf ; \
 in-target sed -i "s#exit 0#sh /root/PnP/bootstrap.sh >> /dev/null 2>\\&1 \\&#g" /etc/rc.local

`
