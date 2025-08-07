<template>
  <DoLoginView v-if="!authenticatedUser" @dologin-success="onDoLoginSuccess" />
  <DashboardView
    v-else
    :user="authenticatedUser"
    @profile-update="onProfileUpdated"
  />
</template>

<script setup>
import { ref } from "vue";
import api from "@/services/api";

import DashboardView from "@/views/DashboardView.vue";
import DoLoginView from "@/views/DoLoginView.vue";

const authenticatedUser = ref(null);

function onDoLoginSuccess(user) {
  authenticatedUser.value = user;
  api.defaults.headers.common["Authorization"] = `Bearer ${user.userId}`;
}

function onProfileUpdated(updatedUser) {
  authenticatedUser.value = updatedUser;
}
</script>
