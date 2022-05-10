# Changelog for unreleased

The following sections list the changes for unreleased.

## Summary

 * Chg #28: Integrate configuration files

## Details

 * Change #28: Integrate configuration files

   We integrated the functionality to support different kinds of configuration files. You can
   find example configurations within the repository. The supported file formats are pretty
   flexible, so far it should work out of the box with `yaml`, `json` and at least `hcl`.

   https://github.com/webhippie/terrastate/issues/28


# Changelog for 1.0.1

The following sections list the changes for 1.0.1.

## Summary

 * Fix #30: Bind flags correctly to variables

## Details

 * Bugfix #30: Bind flags correctly to variables

   We fixed the binding of flags to variables as this had been bound to the root command instead of
   the server command where it belongs to.

   https://github.com/webhippie/terrastate/issues/30


# Changelog for 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #3: Initial release of basic version

## Details

 * Change #3: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/webhippie/terrastate/issues/3


