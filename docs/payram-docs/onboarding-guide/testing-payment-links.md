# Testing Payment Links

***

## Prerequisites

Before generating payment links, ensure the following steps are completed:

* Successfully set up the blockchain node configuration where you will accept payments.
* Complete the wallet management setup so your wallets are ready to receive payments.

***

## Generate payment link&#x20;

{% stepper %}
{% step %}

### Creating payment link

* Go to Payments and expand the section. Select Create Payment Link to generate a new payment link.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/tLQbUkAduKqzIlTlo3Nx/image.png" alt=""><figcaption></figcaption></figure>
{% endstep %}

{% step %}

### How to add a member

* Before you select **Generate Payment Link**, complete these steps:
  1. Add a new member by entering their email, or select an existing member if available.
  2. If this is your first time setting up, the member list is likely empty.
     {% endstep %}

{% step %}

### Add new member

* Select the member email input box. Select Add New Member when the option appears. A pop-up screen opens where you can enter the member’s details.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/7TNIJ6HFo2mIi5sDDr1x/image.png" alt=""><figcaption></figcaption></figure>

* Enter the customer’s email, select the project, and then select Add Member.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/JWqHm7VAlKAik0Qld5pz/image.png" alt=""><figcaption></figcaption></figure>
{% endstep %}

{% step %}

### Enter amount

* After you add the member details, enter the amount to charge that user. Select Generate Payment Link. The system generates a link, which you share with your customer.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/7HzRf41zREm7FXXfAmqW/image.png" alt=""><figcaption></figcaption></figure>

{% endstep %}

{% step %}

### Click on generate payment link

* After you add the member details, enter the amount to charge that user. Select Generate Payment Link. The system generates a link, which you share with your customer.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/NqqoblzV5Ri0aSrs4tsO/image.png" alt=""><figcaption></figcaption></figure>
{% endstep %}

{% step %}

### Payment page

* You'll be redirected to a payment link you can share that to your customer

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/mt5NMtfqeF8745jgxAWU/image.png" alt=""><figcaption></figcaption></figure>

{% hint style="info" %} <mark style="color:$primary;">**Note**</mark><mark style="color:$info;">: If you generate a payment link and the deposit address appears blank, it means the blocks are not being processed. To fix this, restart your PayRam server by running the reset command script.</mark>

👉 [View Restart Command Guide](https://docs.payram.com/script/script-usage)
{% endhint %}
{% endstep %}

{% step %}

### Select preferred coin and network

* When making a payment, customers can select their preferred coin and network.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/3Ol5Pp4kXE4iQquI58Ts/image.png" alt=""><figcaption></figcaption></figure>

* Currently, PayRam supports these coins and networks.
  {% endstep %}

{% step %}

### Payment successful

* After the customer pays, the status updates once the minimum number of onchain confirmations are complete. When confirmed, the system shows Payment Successful.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/Tbkf61bT9GcsRhxtihId/image.png" alt=""><figcaption></figcaption></figure>
{% endstep %}
{% endstepper %}

***

You’ve successfully set up PayRam to accept payments. If you want to integrate the PayRam system into your SaaS, dApp, or any other application, you can do so via the API. Visit the API References section to view the complete integration guide.
