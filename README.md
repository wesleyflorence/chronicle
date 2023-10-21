# chronicle
Tracking Cancer metrics using Notion as a backend.  

## What is this?
A little PWA for my own personal use. I started tracking medication and some health metrics that did not have a home in my current array of apps using Notion. This service was built to make that process easier without moving my existing data. 

### Takeaways
Rather than reach for the familiar, I tried to learn a few new tools. Just to complete the exercise I'll add some commentary here.

#### Go
I usually prefer something more expressive but go has started to grow on me. I can see how the smallness of the language aids in readability and could make you productive quickly. 
I started my project using the Fiber web framework but refactored to Chi because it felt closer to the standard library while still providing some convenience. Both were fine at this tiny scale.
I was not diligent in writing tests during development. I tried to add some for fun and found it a little cumbersome. I think using interfaces more widely would have made mocking less painful. The type system was a joy when writing the app but here I found I was fighting it a bit. I think the experience will change how I structure go code in the future.

#### HTMX
Really nice when it works. It feels dirty to use javascript as an escape hatch for functionality I couldn't get working with HTMX, which I had to do on occasion.

#### Tailwind
I don't write much CSS day to day so this made it very easy and fast to style my page. I'd reach for it again. The standalone cli was great to have so this project could avoid having any package.json.

#### Fly.io
Simplest possible use case here, but deployment with github actions is a breeze.

#### Nix
I'm experimenting with Nix on my machines for a more stable dev environment. I made a simple little flake for this project and it works great. This is a pattern I will use for projects going forward. It solves a similar problem as virtual envs more elegantly. It wouldn't be very useful in a team context though as it would require too much buy in. Docker might be the right tool for this but I haven't gotten that to stick in my workflow for local development.

<img width="600" alt="login_img" src="https://github.com/wesleyflorence/chronicle/assets/28170877/b68ab036-31ef-4282-ab61-b9366d510984">
<img width="600" alt="app_img" src="https://github.com/wesleyflorence/chronicle/assets/28170877/66d23ff7-2b88-4b09-a5db-c63ec470fb52">

