Bugfix: Bind flags correctly to variables

We fixed the binding of flags to variables as this had been bound to the root
command instead of the server command where it belongs to.

https://github.com/webhippie/terrastate/issues/30
