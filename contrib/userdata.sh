#!/bin/bash

cat <<EOM
#cloud-config

write_files:
    - path: /home/core/server.crt
      owner: core:core
      permissions: 0644
      content: |
        $(cat server/server-cert.pem | sed -e "s/^/        /")

    - path: /home/core/server.key
      owner: core:core
      permissions: 0644
      content: |
        $(cat server/server-key.pem | sed -e "s/^/        /")

    - path: /home/core/ca.crt
      owner: core:core
      permissions: 0644
      content: |
        $(cat ca/ca.pem | sed -e "s/^/        /")

coreos:
  update:
    reboot-strategy: off
  units:
    - name: docker.service
      command: restart
      content: |
        [Unit]
        Description=Docker Application Container Engine
        Documentation=http://docs.docker.io
        After=network.target
        [Service]
        ExecStartPre=/bin/mount --make-rprivate /
        # Run docker but don't have docker automatically restart
        # containers. This is a job for systemd and unit files.
        ExecStart=/usr/bin/docker -d -s=btrfs -r=false -D --tlsverify --tlscert=/home/core/server.crt --tlscacert=/home/core/ca.crt --tlskey=/home/core/server.key -H 0.0.0.0:2376

        [Install]
        WantedBy=multi-user.target
EOM
