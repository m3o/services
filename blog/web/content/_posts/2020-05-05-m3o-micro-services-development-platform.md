---
layout:	post
author: Janos Dobronszki
title:	M3O - A cloud platform for microservices
date:	2020-05-05 10:00:00
---

<br>
M3O (the name comes from the word `micro`) is the cloud platform the [Micro](micro.mu) company is releasing soon. In this post we will try to outline what is M3O - why are we doing it and what it will look like.
<br>

## The Why - Cloud computing is too complex

In one sentence, we think the modern cloud ecosystem is too complicated. Gone are the days of launching a successful company on a shared PHP host, or being happy with [restarting ones Ruby server every hour](https://books.google.hu/books?id=ja1KDAAAQBAJ&pg=PA134&lpg=PA134) due to memory leaks, and still becoming literally Twitter.

Somehow between the separate teams working on different parts of the cloud, ease of use for the developer suffers. Or perhaps it's a form of decision paralysis due to the explosion of available technologies. In this fast moving industry one often can't help but follow the technologies provided by FAANGs in a "Nobody got fired for using IMB" spirit.

We are self aware to realise our thinking might run critically close to the [XKCD comic about standards](https://xkcd.com/927/), but we are idealistic (or crazy/bold/ambitious, it's for the dear reader to decide) enough to believe we, a small company can provide a lean approach to other startups of various sizes. Being a small team, we hope to be the David against the Goliaths and provide a unifying, overarching vision for the whole workflow from local to cloud and all steps in between, while adamanatly keeping the Holy Grail of developer happiness in mind.

We did not always want to get into the cloud hosting business. There were some welcome developments over the past long years like lambdas/cloud functions where we almost thought the industry will catch up with and fulfill our wishes and visions so we can focus on building the rest of the stack. However, we realised waiting for years for small parts of our vision to be materialised exceeds our and our community's patience.

Judging the industry by our primary focus on simplicity and developer productivity and happiness (two concepts we believe to be very correlated), we see constant highs and lows. One stride towards a good direction is cancelled by introduction and promotion of overly complex tools, or at least the promotion of overly complex tools to the wrong audience and for the wrong usecase.

Naturally, with millions of us working in this industry, we can't and don't expect everything to go the way we envision - we just enthusiastically wish for it. Our industry is still very young, and we all work hard to make it and the world a better place. That being said, we hope the hard work and decades of collective experience with microservices will result in something that the users will love and put to productive use.

## The What - An opinionated cloud

Micro is and always was, an opinionated framework and ecosystem. Convention over configuration. Easy bootstrapping with zero dependencies locally. Filling in blanks as demands of scaling and resiliency comes up - by switching out implementations of interfaces with more sophisticated ones - that was always the Micro way.

With M3O we plan to keep this approach, working both ways: moving things from local to prod, or prod to local should be a breeze. Micro is particularly well suited for this as an ecosystem built around the use of interfaces and their different implementations. One of our main goals is to make handling multi environments as easy as possible - regardless of where and how the services are being run, managing and interacting with them should feel just like local processes.

We try to make it so when someone learns how to use Micro locally, deploying to M3O will not be further away then a CLI command (assuming an account on M3O exists).
This reuse of the well known and already useful things that are part of the daily workflow will hopefully provide the easiest way for developers to migrate from local to prod or vice versa.

The M3O platform initially will be just a way to create an account and subscription so start using our infrastructure rather than maintaining their own. We plan to introduce more social aspects later like a marketplace for services later, but the current focus is hosting, to enable our users to put their Micro based applications into production without hassles.

For existing users who already have their own live setup, perhaps an easy to set up custom environments (staging, testing or even per engineer environments) will provide value, or if they are perfetly happy with their current setup, maybe they will use M3O for their next project to save bootstrapping time.

Either way, we invite the reader to Sign Up: we will have a free tier for open source code so M3O will be easy to test drive without strings attached.

Cheers

The Micro Team