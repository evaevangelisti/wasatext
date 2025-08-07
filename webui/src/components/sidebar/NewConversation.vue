<template>
  <div class="new-conversation">
    <div class="query">
      <svg viewBox="0 0 24 24" fill="none" class="query-icon">
        <path
          fill-rule="evenodd"
          clip-rule="evenodd"
          d="M15 10.5C15 12.9853 12.9853 15 10.5 15C8.01472 15 6 12.9853 6 10.5C6 8.01472 8.01472 6 10.5 6C12.9853 6 15 8.01472 15 10.5ZM14.1793 15.2399C13.1632 16.0297 11.8865 16.5 10.5 16.5C7.18629 16.5 4.5 13.8137 4.5 10.5C4.5 7.18629 7.18629 4.5 10.5 4.5C13.8137 4.5 16.5 7.18629 16.5 10.5C16.5 11.8865 16.0297 13.1632 15.2399 14.1792L20.0304 18.9697L18.9697 20.0303L14.1793 15.2399Z"
          fill="var(--color-tertiary)"
        />
      </svg>
      <input
        v-model="query"
        type="text"
        placeholder="Search"
        class="query-input"
      >
    </div>
    <button class="new-group-button" @click="emit('show-new-group')">
      <div class="new-group-icon">
        <svg width="24" height="24" viewBox="0 0 1920 1920" fill="none">
          <path
            d="M866.332 213v653.332H213v186.666h653.332v653.332h186.666v-653.332h653.332V866.332h-653.332V213z"
            fill-rule="evenodd"
            fill="var(--color-secondary)"
          />
        </svg>
      </div>
      <span class="text-body">New group</span>
    </button>
    <ul class="users">
      <li v-for="usr in results" :key="usr.userId" class="user">
        <button
          class="user__button"
          :disabled="creating"
          @click="createPrivateConversation(usr.userId)"
        >
          <img
            :src="
              usr?.profilePicture
                ? backendBaseUrl + usr.profilePicture
                : defaultProfilePicture
            "
            alt="Profile picture"
            class="user__picture"
          >
          <span class="text-body">{{ usr.username }}</span>
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import api from "@/services/api";
import { backendBaseUrl } from "@/services/baseUrl";

import defaultProfilePicture from "@/assets/default-profile-picture.jpg";

const props = defineProps({
  user: Object,
});

const emit = defineEmits(["active-conversation", "show-new-group"]);

const query = ref("");
const results = ref([]);
const creating = ref(false);

watch(query, async (newQuery) => {
  if (!newQuery.trim()) {
    results.value = [];
    return;
  }

  try {
    const response = await api.get("/users", {
      params: { q: newQuery },
    });

    results.value = response.data;
  } catch (e) {
    results.value = [];
  }
});

async function createPrivateConversation(otherUserId) {
  if (creating.value) return;
  creating.value = true;

  try {
    const convResponse = await api.get("/conversations");
    const conversations = convResponse.data;

    const existing = conversations.find(c =>
      c.type === "private" &&
      c.participants &&
      c.participants.some(u => u.userId === otherUserId)
    );

    if (existing) {
      emit("active-conversation", existing);
    } else {
      const response = await api.post("/conversations", {
        type: "private",
        userId: otherUserId,
      });

      emit("active-conversation", response.data);
    }
  } catch (e) {
    if (e.response && e.response.status === 409) {
      emit("active-conversation", e.response.data);
    }
  } finally {
    creating.value = false;
  }
}
</script>

<style>
.new-conversation {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  padding: 1rem 1.5rem;
}

.query {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 100px;
  padding: 0.5rem 0.875rem;
  background-color: var(--color-quaternary);
  transition: border 0.1s;
}

.query-input {
  border: none;
  background-color: inherit;
  font-size: 0.95rem;
  color: var(--color-secondary);
}

.query-icon {
  width: 24px;
  height: 24px;
}

.new-group-button {
  display: flex;
  align-items: center;
  gap: 1rem;
  border: none;
  border-radius: 12px;
  width: 100%;
  padding: 0.5rem 0.75rem;
  background-color: inherit;
  color: var(--color-secondary);
  transition: background-color 0.1s;
}

.new-group-button:hover,
.new-group-button:focus {
  outline: none;
  background-color: var(--color-quaternary);
}

.new-group-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border: 2px transparent;
  border-radius: 50%;
  object-fit: cover;
  background-color: var(--color-primary);
}

.users {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  width: 100%;
}

.users {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  width: 100%;
}

.user {
  display: flex;
  width: 100%;
}

.user__button {
  display: flex;
  align-items: center;
  gap: 1rem;
  border: none;
  border-radius: 12px;
  width: 100%;
  padding: 0.5rem 0.75rem;
  background-color: inherit;
  color: var(--color-secondary);
  transition: background-color 0.1s;
}

.user__button:hover,
.user__button:focus {
  outline: none;
  background-color: var(--color-quaternary);
}

.user__picture {
  width: 48px;
  height: 48px;
  border: 2px transparent;
  border-radius: 50%;
  object-fit: cover;
}
</style>
