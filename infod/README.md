infod
-----

Simple server that presents information about the local system
over HTTP, in a simple to consume fashion.


Usage
-----

Compile the server:

    ~$ make

Launch it:

    ~$ ./infod

Now query via curl:

    ~$ curl http://localhost:8000/
    {"FQDN":"shelob.(none)",
     "Interfaces":["lo","eth0","vpn","teredo"],
     "IPv4":["192.168.10.64","10.0.0.200"],
     "IPv6":["fe80::a62:66ff:fe28:6cbf","2001:0:53aa:64c:2cbd:3bde:cafe:beef","fe80::ffff:ffff:ffff"]
    }
    

This dump showed all the information.  Perhaps you only cared about the hostname:

    ~$ curl http://localhost:8000/FQDN

Or just the interfaces?

    ~$ curl http://localhost:8000/Interfaces

The intention is that every key can be queried individually, and dynamically.



Patches
-------

Add in things like "free", "`dpkg --list | grep ^ii | awk '{print $2}'`" and I'll accept them.

I guess I should parse `lsb_release` or `/etc/os-release` at the very least.

Steve
-- 
