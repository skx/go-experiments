plugins/
-------

This directory contains some simple code for testing the plugin
interface included in golang 1.8 beta1, or higher.

The basic idea is that we can externalise code into shared libraries
(read "plugins"), which a driver can load dynamically.
