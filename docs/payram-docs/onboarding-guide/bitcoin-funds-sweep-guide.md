# Bitcoin Funds Sweep Guide

{% stepper %}
{% step %}

### Setup PayRam Merchant mobile app

* Navigate to Settings and click on it. Then, go to the Accounts tab. You will see a QR code labeled Connect to PayRam Mobile App. This QR code is used to sync the app with your PayRam server. Download the PayRam mobile app from the app store separately, then scan this QR code to complete the sync.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/ZCKMh3myTAZVnAFH1wJZ/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/AUPHIhcD52YDNbJaIqdB/image.png" alt=""><figcaption></figcaption></figure>

{% endstep %}

{% step %}

### Login to the PayRam app

* Scan the QR code from your phone. Once scanned, you will be prompted to log in. Use the same root login credentials that you used when setting up the PayRam server web application.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/tpxyBgWZQOz8IFDnelB0/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/0Rjhsn2WfpCp3pIls8iB/image.png" alt=""><figcaption></figcaption></figure>

{% endstep %}

{% step %}

### Set a passcode

Once the above step is completed, you will be prompted to set up a passcode. This passcode will be used to secure access to the PayRam mobile app on your device.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/oOk5xT1CBpmExfO4CBoo/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/2Q2DkK8l6dPuzWODaDug/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/7DifM0HvVyFGYsROF4U5/image.png" alt=""><figcaption></figcaption></figure>

{% endstep %}

{% step %}

### Link wallet

* Now you will see the wallet that you attached during the BTC wallet configuration. You need to select this exact same account so it can be linked here. Click on the link icon, then click on the Link Wallet button. You will be asked to enter the seed phrase, so make sure you have the exact same seed phrase for the account you attached earlier when configuring the BTC wallet. After entering the seed phrase, click Link Wallet to complete the process.

{% hint style="info" %} <mark style="color:$primary;">**Note**</mark><mark style="color:$info;">: The seed phrase will never be stored on any server. It is kept only on your local device storage and protected with encryption.</mark>
{% endhint %}

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/m5oPfKHWy3ljeS1EcOkI/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/ycVXGmEjl3PWP93rk584/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/xyvB1WzbyZ2pqbinsEhC/image.png" alt=""><figcaption></figcaption></figure>

{% endstep %}

{% step %}

### Signing requests tab

* Once you have successfully linked the BTC wallet, you are ready to sweep funds from your BTC deposit wallets to your cold wallet. Go to the Signing Requests tab, where you will see the funds that are ready to be swept. This section will list the deposit addresses where you have received payments, allowing you to sweep them into your cold wallet.
  {% endstep %}

{% step %}

### Signing request rabs sections

In the Signing Requests tab, you will see two sub-tabs

* Pending&#x20;
* In Progress

{% tabs %}
{% tab title="Pending" %}

### **Pending**

* Shows the list of deposit addresses that have received funds and are ready to be swept into your cold wallet.
* These addresses are grouped together in batches for sweeping.
* Clicking on a batch allows you to sweep all funds from the deposit addresses in that batch to your cold wallet.
* For example, if a batch contains 1,400 deposit addresses, all of them have funds ready to sweep.
* If you see two batches, that means there are 2,800 deposit addresses ready to be swept.
* Simply approve and verify these sweeps to move the funds.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/5daGgPr1OL80rtyS8uqb/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/pYTByW8iOWcIxYrm22eT/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/i4IuquIDnSYUepXnphNV/image.png" alt=""><figcaption></figcaption></figure>

{% endtab %}

{% tab title="In progress" %}

### **In progress**

* Shows batches of deposit addresses where sweeping has already started.
* Displays the status of each sweep, including whether all funds have been successfully transferred to your cold wallet.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/eF3uXTGIFYWIG04MpjJK/image.png" alt=""><figcaption></figcaption></figure>
{% endtab %}
{% endtabs %}

{% endstep %}

{% step %}

### Sweeps completed

* Once the sweeps are fully completed, you can check the transaction status by moving from the In Progress tab to the History tab.

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/tHD6nug1GCettMK1dOgs/image.png" alt=""><figcaption></figcaption></figure>

<figure><img src="https://content.gitbook.com/content/8ZOaJUUrFdpVW6K9dkiZ/blobs/xVWZV2O4jlO9YCTW2VWk/image.png" alt=""><figcaption></figcaption></figure>
{% endstep %}
{% endstepper %}
