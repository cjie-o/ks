const { readFileSync } = require("fs");
const { NetCDFReader } = require("netcdfjs");

// http://www.unidata.ucar.edu/software/netcdf/examples/files.html
const data = readFileSync("cru_ts4.05.1901.2020.pre.dat_118-126E_38-44N.nc");

var reader = new NetCDFReader(data); // read the header
reader.getDataVariable("wmoId"); // go to offset and read it
a = reader.attributeExists("pre")
console.log(a)