<template>
    <div class="new-group">
        <div class="new-group__header">
            <button class="back-button" @click="emit('back')">
                <svg viewBox="0 0 24 24" fill="none" class="back-button-icon">
                    <path
                        fill-rule="evenodd"
                        clip-rule="evenodd"
                        d="M11.7071 4.29289C12.0976 4.68342 12.0976 5.31658 11.7071 5.70711L6.41421 11H20C20.5523 11 21 11.4477 21 12C21 12.5523 20.5523 13 20 13H6.41421L11.7071 18.2929C12.0976 18.6834 12.0976 19.3166 11.7071 19.7071C11.3166 20.0976 10.6834 20.0976 10.2929 19.7071L3.29289 12.7071C3.10536 12.5196 3 12.2652 3 12C3 11.7348 3.10536 11.4804 3.29289 11.2929L10.2929 4.29289C10.6834 3.90237 11.3166 3.90237 11.7071 4.29289Z"
                        fill="var(--color-secondary)"
                    />
                </svg>
            </button>
            <span class="text-body">New group</span>
        </div>
        <div class="new-group__content">
            <input
                id="new-group__name"
                type="text"
                placeholder="Group name"
                ref="inputRef"
                v-model="groupName"
                :class="{ 'new-group__name--error': groupNameError }"
            />
            <button class="create-group-button" @click="createGroup">
                <svg viewBox="0 0 24 24" fill="none" class="create-group-icon">
                    <path
                        d="M4 12.6111L8.92308 17.5L20 6.5"
                        stroke="var(--color-secondary)"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import api from "@/services/api";

const emit = defineEmits(["back", "activeConversation"]);

const inputRef = ref(null);
const groupName = ref("");
const groupNameError = ref(false);

onMounted(() => {
    inputRef.value?.focus();
});

async function createGroup() {
    if (!groupName.value.trim()) {
        groupNameError.value = true;
        return;
    }

    groupNameError.value = false;

    try {
        const response = await api.post("/conversations", {
            type: "group",
            name: groupName.value,
            members: [],
        });

        emit("activeConversation", response.data);
    } catch (e) {
        console.error(e);
        groupNameError.value = true;
    }
}
</script>

<style scoped>
.new-group {
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 3rem;
    padding: 0rem 1.5rem;
}

.new-group__header {
    align-self: flex-start;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
}

.back-button {
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 50%;
    padding: 0.3rem;
    background-color: inherit;
}

.back-button:hover {
    background-color: #252526;
    transition: background 0.1s;
}

.back-button-icon {
    width: 20px;
    height: 20px;
}

.new-group__content {
    width: 100%;
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 2.5rem;
    padding: 0rem 1.5rem;
}

#new-group__name {
    border: none;
    border-bottom: 2px solid var(--color-tertiary);
    width: 100%;
    padding: 0.5rem 0rem;
    background-color: inherit;
    font-size: 1rem;
    color: var(--color-secondary);
}

#new-group__name:focus-within {
    border-bottom-color: var(--color-primary);
}

.new-group__name--error {
    border-bottom-color: var(--color-error) !important;
}

.create-group-button {
    align-self: flex-end;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 50%;
    padding: 0.5rem;
    background-color: var(--color-primary);
}

.create-group-icon {
    width: 36px;
    height: 36px;
}

.create-group-button:hover {
    filter: brightness(1.05);
}
</style>
