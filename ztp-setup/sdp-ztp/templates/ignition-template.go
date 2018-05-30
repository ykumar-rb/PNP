package templates

//template for ignition
var IgnitionTmlp = `---
systemd:
  units:
    - name: installer.service
      enable: true
      contents: |
        [Service]
        Type=oneshot
        ExecStart=/usr/bin/echo Hello > /tmp/1.txt
        [Install]
        WantedBy=multi-user.target

`
