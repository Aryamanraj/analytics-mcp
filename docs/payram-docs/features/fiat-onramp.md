# Fiat Onramp

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2Fj9r37tuqq1mslKbOZI3F%2Fpayram-feature-fiat-onramp.png?alt=media&#x26;token=605f97c8-5c5d-4d1d-b3e3-90760d369982" alt=""><figcaption></figcaption></figure>

Fiat Onramp allows merchants to accept customer payments in fiat currency while receiving settlements in crypto through PayRam. It removes the need for external exchanges or manual conversions, letting businesses expand their customer base and simplify checkout experiences.

***

### Why it Matters

* Lets merchants reach customers who prefer paying in fiat.
* Enables seamless acceptance of cards, wallets, bank transfer, and local payment methods.
* Automatically settles payments in crypto, keeps all settlements non-custodial.
* Reduces conversion friction and simplifies operations.

***

### Supported onramp methods

### Third-party Onramp Partners

This method connects merchants to regulated fiat-to-crypto providers, enabling customers to complete transactions using familiar fiat payment options like card, bank transfer, and more.

#### How it works for merchants

* To enable the onramp feature in PayRam, merchants must first complete KYB with the onramp partner.
* After KYB approval, the merchant needs to add the API key provided by the onramp partner into the PayRam dashboard.
* PayRam will auto-detect the integration and update the available checkout options without any extra configuration.

#### Commercials and fees

* PayRam does not apply any additional fees or markups on onramp transactions.
* All onramp related commercials are directly between the merchant and the onramp provider.

#### **Checkout experience**

* All onramp transaction amounts are directly deposited into the merchant’s deposit wallet (unique address created for each customer).

#### Supported partners

* PayRam currently integrates with **TransFi**, more onramp partners will be arriving soon!
* TransFi supports **100+ countries** and **40+ fiat currencies**, allowing global customers to purchase crypto seamlessly.

### PayRam Payments App \[Coming soon]

This method will allow merchants to accept fiat payments through PayRam’s integrated payments layer, offering a faster setup and a smoother customer experience.

#### How it works for merchants:

* Merchants do not need to complete KYC/KYB to enable this Onramp method.
* Activation will be available directly inside the PayRam Dashboard once the feature goes live.
* Merchants will have the option to sponsor gas fees for customers, reducing friction and improving conversion.

#### How it works for customers:

* Customers will still be required to complete KYC verification once before making a purchase.
* After verification, customers can pay using supported fiat methods with minimal steps.

***

## How to Enable Third-party Fiat Onramp Partner

{% stepper %}
{% step %}

### Navigate to settings

* Log in to your PayRam Dashboard and go to the Settings section.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2F1U2ZAxVVmN3ilrl7MA5F%2Fimage.png?alt=media&#x26;token=203c05e8-bb21-40eb-8f69-19628e9adde3" alt=""><figcaption></figcaption></figure>

* Now click on the Payment Channels option.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2FVbg9kJy4PGaChyZOL17o%2Fimage.png?alt=media&#x26;token=c3fc2898-e20a-4df2-af2e-5ac14ae6fe95" alt=""><figcaption></figcaption></figure>

{% endstep %}

{% step %}

### Activate TransFi

* Click on the Activate button beside *TransFi* to enable this payment channel.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2FFA2o81t0B57BJSyaOvoW%2Fimage.png?alt=media&#x26;token=5559b993-b70f-4287-860a-f8f387bd1288" alt=""><figcaption></figcaption></figure>

* Get your TransFi API key from their dashboard. You must complete KYB/KYC with TransFi in order to procure the API key.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2FU0S25ZH1PxPwVzymGJk8%2Fimage.png?alt=media&#x26;token=b0b7dd90-a8f3-4969-a4de-7938ad25fd36" alt=""><figcaption></figcaption></figure>

* After you receive the TransFi API key, paste it into the API Key field and click Activate. Once activated, you can start accepting customer payments using the fiat on-ramp feature.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2F8DlEvzcaAATp9lRhpBZu%2Fimage.png?alt=media&#x26;token=fb956495-b3f3-4680-8e80-ec254ac8e39e" alt=""><figcaption></figcaption></figure>
{% endstep %}

{% step %}

### Accept payments

* Go to the Payments menu in the sidebar, click the dropdown, and then select Create Payment Link.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2F6cSYVcKOmz9BbzeiBKPo%2Fimage.png?alt=media&#x26;token=4da1c5a6-bbff-45e7-af93-77be2cd8f99f" alt=""><figcaption></figcaption></figure>

* Create a payment link by entering the customer’s email and the required amount, then click Generate Payment Link.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2FMIUSHUX1pGpQGy9bh7kz%2Fimage.png?alt=media&#x26;token=08e65d2f-5d55-4c7d-bb08-688d795be163" alt=""><figcaption></figcaption></figure>
{% endstep %}

{% step %}

### Pay using TransFi Widget

* Your customers will now see the TransFi payment option on the payment page. They just need to click on TransFi to use it.

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2FoJQIMYcyJwKBCAF9Ggf8%2Fimage.png?alt=media&#x26;token=669fd4ab-a474-42aa-99e8-eaace4acd847" alt=""><figcaption></figcaption></figure>

* Customers can pay using their credit or debit cards through the TransFi widget, and the crypto will be deposited directly into the merchant’s deposit wallet (unique address created for each customer).

<figure><img src="https://2353895776-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F8ZOaJUUrFdpVW6K9dkiZ%2Fuploads%2Fb1PCvPK4k3s3K7ES4kAL%2Fimage.png?alt=media&#x26;token=3a7f5ca5-4125-4b1c-a415-47c9f5ce1d67" alt=""><figcaption></figcaption></figure>
{% endstep %}
{% endstepper %}
