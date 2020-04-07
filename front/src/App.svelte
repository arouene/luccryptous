<script>
    import { Tabs, Tab, TabLabel, TabContent } from './Tabs';

    let ciphertext = "";
    let plaintext = "";
    let password = "";
    let uuid = "";

    async function cryptPlaintext(event) {
        ciphertext = "Waiting...";
        const res = await fetch(`http://localhost:3000/api/crypt`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({secret: event.target.value})
        });
        const pass = await res.json();

        if (res.ok) {
            ciphertext = pass.secret;
        } else {
            ciphertext = "Error...";
        }
    }

    async function genPass() {
        password = "Waiting...";
        const res = await fetch(`http://localhost:3000/api/pass`);
        const pass = await res.json();

        if (res.ok) {
            password = pass.secret;
        } else {
            password = "Error...";
        }
    }

    async function genUUID() {
        uuid = "Waiting...";
        const res = await fetch(`http://localhost:3000/api/uuid`);
        const pass = await res.json();

        if (res.ok) {
            uuid = pass.secret;
        } else {
            uuid = "Error...";
        }
    }

    function switchTab(event) {
        switch(event.detail.tab) {
        case 2:
            genPass();
            break;
        case 3:
            genUUID();
            break;
        }
    }
</script>

<style>
    :global(body) {
        background: linear-gradient(45deg,#ff8737 29.25%,#f60 100%) no-repeat fixed;
        background-color: #f90;
    }

	h1 {
		color: white;
		text-transform: uppercase;
        text-align: center;
		font-size: 4em;
		font-weight: 100;
	}

    input {
        width: 100%;
        font-size: 0.95em;
    }
</style>

<h1>Luccryptous</h1>

<Tabs>
  <Tab on:switch={switchTab}>
    <TabLabel>Chiffrer</TabLabel>
    <TabContent>
      <label>
        <h3>Plaintext :</h3>
        <input on:input={cryptPlaintext}>
      </label>
      <label>
        <h3>Ciphertext :</h3>
        <input readonly value={ciphertext} on:click={navigator.clipboard.writeText(ciphertext)}>
      </label>
    </TabContent>
  </Tab>
  <Tab on:switch={switchTab}>
    <TabLabel>Password</TabLabel>
    <TabContent>
      <label>
        <h3>Ciphertext :</h3>
        <input readonly value={password} on:click={navigator.clipboard.writeText(password)}>
      </label>
      <button on:click={genPass}>Regénérer</button>
    </TabContent>
  </Tab>
  <Tab on:switch={switchTab}>
    <TabLabel>UUID</TabLabel>
    <TabContent>
      <label>
        <h3>Ciphertext :</h3>
        <input readonly value={uuid} on:click={navigator.clipboard.writeText(uuid)}>
      </label>
      <button on:click={genUUID}>Regénérer</button>
    </TabContent>
  </Tab>
</Tabs>
