# Testpoint

**Testpoint** is a simple CLI testing tool written in Go that can help you ensure that your endpoints work 
as expected after major refactoring or migration to a new version of your favourite framework.

The tool has two main features:
* Sending prepared requests to the endpoints you want to test and collecting the responses.
* Comparing the collected responses and generating a report.

## Motivation

I decided to implement this tool after my team and I faced the exact same issue several times: 
we made major changes to an application, which could potentially break everything, 
and we couldn't trust the automated tests that were already written because they didn't cover all the functionality.

One of the common solutions that we had was writing a Python script that sends prepared requests 
(we usually took them from production access logs) to both the verified version of the app and the new one. 
The script collected the responses and then compared them to ensure that there are no differences 
and the new version of the app is safe to release and deploy in production.

We often performed the same kind of testing, in addition to automated tests, when we needed to rewrite a legacy application. 
This approach helped us catch many bugs that were quite difficult to detect otherwise.

After writing multiple Python scripts like that, I realised that we were wasting our time by doing the same work over and over again.
Moreover, the scripts were quite slow (it was Python, after all), and we had to waste even more time waiting for them to finish.
That's when I came up with the idea for this tool.

## When you should use it

Testpoint can be useful in the following cases:
* You've made significant changes that don't alter the behaviour of the endpoints in question, 
  and you need to test that they still work as expected.
* You've rewritten a legacy application, and you need to ensure that the endpoints respond in exactly the same way.

> Note that not every REST endpoint is suitable for this kind of testing. If you want to test an endpoint, 
> make sure that it's **idempotent** and **consistent**, i.e., it produces the same responses regardless of the order or number of requests.
