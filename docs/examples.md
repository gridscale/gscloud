# Examples

## Create a server object with storage

    $ gscloud server create \
        --name test-1 \
        --cores=1 \
        --mem=1 \
        --with-template="CentOS 8 (x86_64)" \
        --hostname test-1
    Server created: f75a03d9-dcd8-4ffe-9d22-38481fc0f4cb
    Storage created: ad1ec283-99a8-4d0a-a219-4f61e7b4a654
    Password: l548F37u1c:^

## Pull an IP address and assign to a server

    $ gscloud ip add -6
    IP added: 2001:db8:0:1::1c8

    $ gscloud ip assign --to=f75a03d9-dcd8-4ffe-9d22-38481fc0f4cb 2001:db8:0:1::1c8

## List servers

    $ gscloud server ls
    ID                                    NAME    CORE  MEM  CHANGED                    POWER
    f75a03d9-dcd8-4ffe-9d22-38481fc0f4cb  test-1  1     1    2021-02-19T20:07:59+01:00  off

## List storages

    $ gscloud storage ls
    ID                                    NAME    CAPACITY  CHANGED                    STATUS
    ad1ec283-99a8-4d0a-a219-4f61e7b4a654  test-1  10        2021-02-19T20:08:04+01:00  active

## List IPs

    % gscloud ip ls
    IP                  ASSIGNED  FAILOVER  FAMILY  REVERSE DNS                                      ID
    2001:db8:0:1::1c8  assigned  no        v6      static-2001:db8:0:1::1c8.gridserver.io      3917cff6-0f7d-408c-8918-dfe1eae498ea

## Remove a server along with assigned IPs and storages

    $ gscloud server rm -i f75a03d9-dcd8-4ffe-9d22-38481fc0f4cb
    ID                                    TYPE          NAME
    f75a03d9-dcd8-4ffe-9d22-38481fc0f4cb  Server        test-1
    ad1ec283-99a8-4d0a-a219-4f61e7b4a654  Storage       test-1
    816b9584-5a68-4cd3-a77c-a4fcc28ed717  IPv6 address  2001:db8:0:1::1c8
    This can destroy your data. Re-run with --force to remove above objects

And if you are really sure what you do:

    $ % gscloud server rm -i -f 2090776c-f428-40b4-a97b-f61081516eb2
    ID                                    TYPE          NAME
    2090776c-f428-40b4-a97b-f61081516eb2  Server        test-1
    885537a6-a836-47bd-ad1a-0d872e2a2317  Storage       test-1
    3917cff6-0f7d-408c-8918-dfe1eae498ea  IPv6 address  2001:db8:0:1::1c8
    Removed 2090776c-f428-40b4-a97b-f61081516eb2
    Removed 885537a6-a836-47bd-ad1a-0d872e2a2317
    Removed 3917cff6-0f7d-408c-8918-dfe1eae498ea
