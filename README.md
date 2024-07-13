# Testpoint

**Testpoint** is a simple CLI testing tool written in Go that can help ensure that your endpoints work as expected after
major refactoring or migration to a new version of your favourite framework.

The tool has two main features:

* Sending prepared requests to the endpoints you want to test and collecting the responses.
* Comparing the collected responses and generating a report.

## Motivation

I decided to implement this tool after my team and I faced the exact same issue several times: we made major changes to
an application, which could potentially break everything, and we couldn't trust the automated tests that were already
written because they didn't cover the full range of the application’s functionality.

One of the common solutions that we had was writing a Python script that sends prepared requests (we usually took them
from production access logs) to both the verified version of the app and the new one. The script collected the responses
and then compared them to ensure that there were no differences and the new version of the app is safe to release and
deploy in production.

We often performed the same kind of testing, in addition to automated tests, when we needed to rewrite a legacy
application. This approach helped us catch many bugs that would have been quite difficult to detect otherwise.

After writing multiple Python scripts like that, I realised that we were wasting our time by doing the same work over
and over again. Moreover, the scripts were quite slow (it was Python, after all), and we had to waste even more time
waiting for them to finish. That's why I decided to create Testpoint.

## When you should use it

Testpoint can be useful in the following scenarios:

* You've made significant changes that don't alter the behaviour of the endpoints in question, and you need to test that
  they still work as expected.
* You've rewritten a legacy application, and you need to ensure that the endpoints respond in exactly the same way.

Note that not every REST endpoint is suitable for this kind of testing. If you want to test an endpoint, make sure that
it's **idempotent** and **consistent**, i.e., it produces the same responses regardless of the order or number of
requests.

## Sending requests

Let's assume you've already prepared a CSV file with requests and named it `requests.csv`. It might look something like
this:

```
url
https://test.com/api/v1/suggestions?prefix=at
https://test.com/api/v1/suggestions?prefix=ca
https://test.com/api/v1/suggestions?prefix=to
https://test.com/api/v1/suggestions?prefix=ta
https://test.com/api/v1/suggestions?prefix=ru
https://test.com/api/v1/suggestions?prefix=ga
https://test.com/api/v1/suggestions?prefix=tr
https://test.com/api/v1/suggestions?prefix=ch
```

Use the following command to send the requests to the desired hosts and collect the responses:

```shell
testpoint send ./requests.csv http://localhost:8083 http://localhost:8084
```

The `send` command takes several arguments: the first one is a file or a directory with your requests, and the following
are the URLs of the applications you want to test. (Note that you can specify any number of URLs; it's not strictly
necessary to have two as shown in the example)

When the processing is completed, you'll find the output files with collected responses in the same directory where the
command was executed. Typically, the names of the output files are based on the names of the given URLs;
for example, `http-localhost-8083.csv` and `http-localhost-8084.csv`.

### Additional request data

You can also specify request methods, headers, and bodies in your CSV file:

```
url,method,headers,body
https://test.com/api/v1/suggestions?prefix=at,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=ca,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=to,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=ta,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=ru,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=ga,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=tr,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
https://test.com/api/v1/suggestions?prefix=ch,GET,"{""headerField"":""test""}","{""bodyField"":""test""}"
```

The order of the columns doesn't matter; the values will be found by column names.

Note that the command uses the `GET` method by default if it's not specified otherwise in the file.

### URL substitution

As you might have noticed, the requests from the CSV file already include the host, which is `https://test.com`.
However, this is not a problem because the `send` command knows which part of the URL needs to be replaced before
sending the request.

That being said, it's perfectly fine to specify your request URLs in the following format as well:

```
url
/api/v1/suggestions?prefix=at
/api/v1/suggestions?prefix=ca
/api/v1/suggestions?prefix=to
/api/v1/suggestions?prefix=ta
/api/v1/suggestions?prefix=ru
/api/v1/suggestions?prefix=ga
/api/v1/suggestions?prefix=tr
/api/v1/suggestions?prefix=ch
```

Moreover, you can also replace the URL path if it's necessary.
To do that, you just need to include it to the URL when you run the command:

```shell
testpoint send ./requests.csv http://localhost:8083/new-endpoint
```

So, for instance, instead of `http://localhost:8083/api/v1/suggestions?prefix=at`,
the request will be sent to `http://localhost:8083/new-endpoint?prefix=at`.

### Workers

By default, the `send` command uses only one thread to send requests. However, if you have a lot of input data, the
execution might take a while. To speed it up, you might want to increase the number of workers (you can think of them as
threads) using the `--workers` or just `-w` flag:

```shell
testpoint -w 8 send ./requests.csv http://localhost:8083
```

### Custom request transformation

The default request transformation is usually sufficient for most cases. However, if your request data is arranged
differently in the CSV file or if you need to make specific changes to your requests before sending them, you can always
write your own custom transformation using JavaScript.

As an example, let's create a new transformation that will read our custom columns from the CSV file.
Here's the input file:

```
path,prefix,method
/api/v1/suggestions,at,GET
/api/v1/suggestions,ca,GET
/api/v1/suggestions,to,GET
/api/v1/suggestions,ta,GET
/api/v1/suggestions,ru,GET
/api/v1/suggestions,ga,GET
/api/v1/suggestions,tr,GET
/api/v1/suggestions,ch,GET
```

Next, we need to create a new JavaScript file with the `transform` function that takes two arguments:
the host, which we specify when running the command, and the CSV record from the input file:

```javascript
function transform(host, record) {
    return {
        url: host + record["path"] + "?prefix=" + record["prefix"],
        method: record["method"],
    };
}
```

The returning value is an object containing `url`, `method`, `headers`, and `body`. If some properties are not needed,
you can leave them out.

Finally, you can run the `send` command with the `--transformation` or simply `-t` flag to specify the new
transformation:

```shell
testpoint send -t transformation.js ./requests.csv http://localhost:8083 http://localhost:8084
```

Note that if you implement your own custom transformation, you need to take care of the URL substitution yourself
because it's a feature of **the default transformation**.

## Comparing responses

After collecting the responses, you might want to compare them to see if there are any differences. To do that, run
the `compare` command with the two csv files as arguments:

```shell
testpoint compare ./http-localhost-8083.csv ./http-localhost-8084.csv
```

If there are any differences in responses, the mismatches will be printed in your terminal like this:
![Mismatch Example](https://i.imgur.com/SnmbUvh.png)

As you can see in the screenshot, there are a few differences between the two responses: the JSON object with `id` 42
has appeared, and the object with `id` 45 is no longer there.

If you want to collect all the mismatched responses into a file, you can add the `--csv-report` flag when you run the command,
specifying the name of the output CSV file:

```shell
testpoint compare --csv-report ./report.csv ./http-localhost-8083.csv ./http-localhost-8084.csv
```

### Custom comparator

TBD
