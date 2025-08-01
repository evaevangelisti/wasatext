<template>
    <div class="sidebar" @photo-updated="refreshUser">
        <SidebarHeader
            :user="user"
            :mode="mode"
            @showConversationList="mode = 'conversationList'"
            @showNewConversation="mode = 'newConversation'"
            @showProfile="mode = 'profile'"
        />
        <ConversationList v-if="mode === 'conversationList'" />
        <NewConversation v-if="mode === 'newConversation'" />
        <Profile
            :user="user"
            v-if="mode === 'profile'"
            @profile-updated="$emit('profile-updated', $event)"
        />
    </div>
</template>

<script setup>
import { ref } from "vue";

import SidebarHeader from "@/components/Sidebar/SidebarHeader.vue";
import ConversationList from "@/components/Sidebar/ConversationList.vue";
import NewConversation from "@/components/Sidebar/NewConversation.vue";
import Profile from "@/components/Sidebar/Profile.vue";

const props = defineProps({
    user: Object,
});

const emit = defineEmits(["profile-updated"]);

const mode = ref("conversationList");
</script>

<style setup>
.sidebar {
    display: flex;
    flex-direction: column;
    border-right: 1px solid var(--color-tertiary);
    width: 30vw;
}
</style>
