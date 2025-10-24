<template>
  <DoLoginView v-if="!authenticatedUser" @dologin-success="onDoLoginSuccess" />
  <DashboardView
    v-else
    :user="authenticatedUser"
    @profile-update="onProfileUpdated"
    @logout="onLogout"
  />
</template>

<script setup>
import { ref, onMounted } from "vue";
import api from "@/services/api";

import DashboardView from "@/views/DashboardView.vue";
import DoLoginView from "@/views/DoLoginView.vue";

const authenticatedUser = ref(null);

function onDoLoginSuccess(user) {
  authenticatedUser.value = user;
  api.defaults.headers.common["Authorization"] = `Bearer ${user.userId}`;
  localStorage.setItem("authenticatedUser", JSON.stringify(user));
}

function onProfileUpdated(updatedUser) {
  authenticatedUser.value = updatedUser;
  localStorage.setItem("authenticatedUser", JSON.stringify(updatedUser));
}

function onLogout() {
  authenticatedUser.value = null;
  delete api.defaults.headers.common["Authorization"];
  localStorage.removeItem("authenticatedUser");
}

onMounted(() => {
  const userStr = localStorage.getItem("authenticatedUser");
  if (userStr) {
    const user = JSON.parse(userStr);
    authenticatedUser.value = user;
    api.defaults.headers.common["Authorization"] = `Bearer ${user.userId}`;
  }
});
</script>
