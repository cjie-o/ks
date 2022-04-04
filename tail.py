# from pylab import *
import netCDF4

f = netCDF4.MFDataset('./*.nc')
# print variables
f.variables.keys()
atemp = f.variables['air'] # TODO: extract spatial subset

latbounds = [ 40 , 43 ]
lonbounds = [ -96 , -89 ] # degrees east ? 
lats = f.variables['latitude'][:] 
lons = f.variables['longitude'][:]

# latitude lower and upper index
latli = np.argmin( np.abs( lats - latbounds[0] ) )
latui = np.argmin( np.abs( lats - latbounds[1] ) ) 

# longitude lower and upper index
lonli = np.argmin( np.abs( lons - lonbounds[0] ) )
lonui = np.argmin( np.abs( lons - lonbounds[1] ) )  