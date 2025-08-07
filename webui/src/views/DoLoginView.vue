<template>
  <div class="do-login">
    <div class="do-login__header">
      <span class="text-title do-login__title">WASAText</span>
      <span class="text-subtitle do-login__subtitle">Sign in</span>
    </div>
    <form class="do-login__form" @submit.prevent="doLogin">
      <div class="do-login__field">
        <input
          ref="usernameInput"
          v-model="username"
          placeholder="Username"
          autocomplete="off"
          class="do-login__input"
        >
        <span class="text-caption error">{{ error || "\u00A0" }}</span>
      </div>
      <button type="submit" class="submit">Next</button>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import api from "@/services/api";

const username = ref("");
const error = ref("");
const usernameInput = ref(null);

const emit = defineEmits(["dologin-success"]);

async function doLogin() {
  error.value = "";

  try {
    const response = await api.post("/users", { username: username.value });
    emit("dologin-success", response.data);
  } catch (e) {
    error.value = "Invalid username";
  }
}

onMounted(() => {
  if (usernameInput.value) {
    usernameInput.value.focus();
  }
});
</script>

<style scoped>
.do-login {
  display: flex;
  justify-content: center;
  gap: 15rem;
  border-radius: 1.75rem;
  padding: 4rem;
  background-color: #000000;
}

.do-login__header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.do-login__title {
  font-size: 2.25rem;
  color: var(--color-primary);
}

.do-login__subtitle {
  font-size: 1.5rem;
}

.do-login__form {
  display: flex;
  align-items: flex-end;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
}

.do-login__field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.do-login__input {
  border: 1px solid var(--color-tertiary);
  border-radius: 4px;
  padding: 0.75rem 0.875rem;
  background-color: inherit;
  font-size: 1rem;
  color: var(--color-secondary);
  transition: border 0.1s;
}

.do-login__input:focus {
  border-color: var(--color-primary);
}

.submit {
  border: none;
  border-radius: 28px;
  padding: 0.625rem 1.25rem;
  background-color: var(--color-primary);
  font-size: 1rem;
  color: var(--color-secondary);
  transition: filter 0.1s;
}

.submit:hover {
  filter: brightness(1.05);
}
</style>
