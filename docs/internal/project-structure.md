# Project Structure

| Path           	| Usage                                                                                                                                              	|
|----------------	|----------------------------------------------------------------------------------------------------------------------------------------------------	|
| /configuration 	| This package should contain anything related to setting up the application from the config file                                                    	|
| /controller    	| This package should contain anything related to handling, creating and attaching Controller's                                                      	|
| /database      	| This package should only contain setup for the database                                                                                            	|
| /docs          	| This directory contains the documentation for the project. Internal is for developer's to refrence when building.  External for user's.            	|
| /e2e           	| This directory contain's end-to-end test's, E2E tests are written in ruby                                                                          	|
| /misc          	| This package should contain utility function's like Capitalize etc.                                                                                	|
| /model         	| This package should contain struct's, methods and function's used in generating json spec,  handling data and other component-specific function's. 	|