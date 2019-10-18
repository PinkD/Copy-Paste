# Copy Paste

Copy and paste code

## Usage

```bash
cat file | curl -F "content=<-" you_host/new
```

## Library And Resources

`highlight.pack.js` is from [highlight.js](https://github.com/highlightjs/highlight.js)

`highlightjs-line-numbers.min.js` is from [highlightjs-line-numbers.js](https://github.com/wcoder/highlightjs-line-numbers.js/blob/master/src/highlightjs-line-numbers.js)

`favicon.ico` is from Internet

## Benchmark

### post
- new content(HDD)
  - 1m39.4464003s elapsed, 114840 lines sent, 1154.79 op/s
- new content(memory)
  - 13.6784482s elapsed, 114840 lines sent, 8395.69 op/s
- new exists content
  - 6.6683451s elapsed, 114840 lines sent, 17221.66 op/s

## get
- with redis
  - 5.302129s elapsed, 111815 lines got, 21088.70 op/s
- without redis(memory)
  - 10.1018404s elapsed, 111815 lines got, 11068.78 op/s
- without redis(HDD)
  - 10.659549s elapsed, 111815 lines got, 10489.66 op/s


## License

```license
Copyright 2019 PinkD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```