[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/AntoineAugusti/avurnav-api)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/AntoineAugusti/avurnav-api/blob/master/LICENSE.md)
![Swagger Validator](https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/AntoineAugusti/avurnav-api/master/docs.yml.svg)


# AVURNAV API
This library exposes an HTTP API to get navigational warnings for Metropolitan France currently active. Navigational warnings contain information about persons in distress, or objects and events that pose an immediate hazard to navigation. Navigational warnings are called [AVURNAV](https://fr.wikipedia.org/wiki/Avis_urgent_aux_navigateurs) (avis urgent aux navigateurs) in French.

It relies on another library: https://github.com/AntoineAugusti/avurnav.

On top of exposing an HTTP API, it periodically fetches information from the Préfet Maritime websites to get the latest published AVURNAVs.

## API documentation
An HTML documentation of the API is available [here](https://antoineaugusti.github.io/avurnav-api/). The documentation is build thanks to [Spectacle](https://github.com/sourcey/spectacle).

This documentation is build using the Swagger 2.0 specification. With Swagger, you can automatically build HTTP clients for this API. The definition file can be found [here](https://github.com/AntoineAugusti/avurnav-api/blob/master/docs.yml).

## Notice
This software is available under the MIT license and was developed as part of the [Entrepreneur d'Intérêt Général program](https://entrepreneur-interet-general.etalab.gouv.fr) by the French government.
