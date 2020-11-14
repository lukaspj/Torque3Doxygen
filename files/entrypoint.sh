#!/bin/bash
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP
sudo -u app sh -c "$@"