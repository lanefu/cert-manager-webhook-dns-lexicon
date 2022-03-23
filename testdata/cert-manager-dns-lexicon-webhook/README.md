# Solver testdata directory

Running the tests:
------------------

* First you need python 3 installed and the dns-lexicon library (`pip install dns-lexicon`)
* copy apikey.yaml.sample to apikey.yaml
* Update apikey.yaml with the keys that you want to use
* Update the provider in config.json to the correct provider (see [DNS lexicon docs](https://dns-lexicon.readthedocs.io/en/latest/introduction.html#supported-providers)) 
* optionally set sandbox in config.json to true/false as desired

    TEST_ZONE_NAME=mydomain.tld. make test