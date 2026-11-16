<script setup>
import { ref, computed, onMounted } from 'vue'
import { userService } from './services/api'
import UserFormModal from './components/UserFormModal.vue'

const users = ref([])
const searchQuery = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const isFormOpen = ref(false)
const selectedUser = ref(null)
const selectedIds = ref([])

let successTimeout = null

const filteredUsers = computed(() => {
  const term = searchQuery.value.trim().toLowerCase()
  if (!term) return users.value
  return users.value.filter(
    (u) =>
      u.name.toLowerCase().includes(term) ||
      u.email.toLowerCase().includes(term) ||
      u.login.toLowerCase().includes(term)
  )
})

function showSuccess(message) {
  successMessage.value = message
  clearTimeout(successTimeout)
  successTimeout = setTimeout(() => {
    successMessage.value = ''
  }, 3000)
}

function toggleSelectAll(event) {
  if (event.target.checked) {
    selectedIds.value = filteredUsers.value.map((u) => u.id)
  } else {
    selectedIds.value = []
  }
}

async function fetchUsers() {
  loading.value = true
  errorMessage.value = ''
  try {
    users.value = await userService.getAll()
    selectedIds.value = []
  } catch (err) {
    errorMessage.value = err.message
  } finally {
    loading.value = false
  }
}

function openCreateForm() {
  selectedUser.value = null
  successMessage.value = ''
  errorMessage.value = ''
  isFormOpen.value = true
}

function openEditForm(user) {
  selectedUser.value = user
  successMessage.value = ''
  errorMessage.value = ''
  isFormOpen.value = true
}

function handleSaved(message) {
  isFormOpen.value = false
  showSuccess(message)
  fetchUsers()
}

async function handleDeleteSelected() {
  const count = selectedIds.value.length
  if (count === 0) return

  if (!confirm(`Deseja realmente excluir ${count} usuário(s) selecionado(s)?`)) return

  errorMessage.value = ''
  try {
    await Promise.all(selectedIds.value.map((id) => userService.remove(id)))
    selectedIds.value = []
    showSuccess(`${count} usuário(s) excluído(s) com sucesso!`)
    fetchUsers()
  } catch (err) {
    errorMessage.value = err.message
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="min-vh-100 py-5 px-3" style="background-color: #e9eef5;">
    <div class="container" style="max-width: 980px;">
      <!-- abre modal de formulario -->
      <UserFormModal
        v-if="isFormOpen"
        :user-to-edit="selectedUser"
        @saved="handleSaved"
        @close="isFormOpen = false"
      />

      <!-- listagem principal -->
      <div v-else class="card border-0 shadow-sm rounded-4 overflow-hidden">
        <div
          class="card-header bg-white border-bottom py-4 px-4 d-flex flex-wrap justify-content-between align-items-center gap-3"
        >
          <h1 class="h4 fw-bold mb-0 text-dark">User management</h1>

          <div class="flex-grow-1" style="max-width: 280px;">
            <input
              v-model="searchQuery"
              type="text"
              class="form-control form-control-sm bg-light border-0 px-3 py-2 rounded-3"
              placeholder="Buscar por nome, e-mail ou login..."
            />
          </div>

          <div class="d-flex align-items-center gap-2">
            <button
              v-if="selectedIds.length > 0"
              class="btn btn-outline-danger btn-sm px-3 py-2 rounded-3 fw-medium"
              @click="handleDeleteSelected"
            >
              Excluir ({{ selectedIds.length }})
            </button>

            <button
              class="btn btn-primary btn-sm px-3 py-2 rounded-3 fw-medium"
              @click="openCreateForm"
            >
              + Novo usuário
            </button>
          </div>
        </div>

        <!-- sucesso/erro -->
        <div v-if="successMessage || errorMessage" class="px-4 pt-3">
          <div v-if="successMessage" class="alert alert-success py-2 mb-0 rounded-3" role="alert">
            {{ successMessage }}
          </div>
          <div v-if="errorMessage" class="alert alert-danger py-2 mb-0 rounded-3" role="alert">
            {{ errorMessage }}
          </div>
        </div>

        <!-- loading -->
        <div v-if="loading" class="text-center py-5 text-muted">
          <div class="spinner-border spinner-border-sm text-primary me-2" role="status"></div>
          <span>Carregando usuários...</span>
        </div>

        <!-- table -->
        <div v-else-if="filteredUsers.length > 0" class="table-responsive">
          <table class="table table-hover align-middle mb-0">
            <thead class="bg-light">
              <tr class="text-uppercase text-muted" style="font-size: 0.75rem; letter-spacing: 0.04em;">
                <th class="ps-4 py-3" style="width: 44px;">
                  <input
                    type="checkbox"
                    class="form-check-input"
                    :checked="filteredUsers.length > 0 && selectedIds.length === filteredUsers.length"
                    @change="toggleSelectAll"
                  />
                </th>
                <th class="py-3 fw-semibold">Nome</th>
                <th class="py-3 fw-semibold">E-mail</th>
                <th class="py-3 fw-semibold">Login</th>
                <th class="py-3 fw-semibold">Status</th>
                <th class="pe-4 py-3 text-end fw-semibold">Ações</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.id">
                <td class="ps-4">
                  <input
                    type="checkbox"
                    class="form-check-input"
                    :value="user.id"
                    v-model="selectedIds"
                  />
                </td>
                <td class="fw-medium text-dark py-3">{{ user.name }}</td>
                <td class="text-secondary">{{ user.email }}</td>
                <td class="text-secondary">{{ user.login }}</td>
                <td>
                  <span class="d-inline-flex align-items-center gap-2 small fw-medium">
                    <span
                      class="rounded-circle d-inline-block"
                      :class="user.active ? 'bg-success' : 'bg-danger'"
                      style="width: 8px; height: 8px;"
                    ></span>
                    <span :class="user.active ? 'text-success' : 'text-danger'">
                      {{ user.active ? 'Ativo' : 'Inativo' }}
                    </span>
                  </span>
                </td>
                <td class="pe-4 text-end">
                  <button
                    class="btn btn-light btn-sm border px-3 rounded-3 text-dark fw-medium"
                    @click="openEditForm(user)"
                  >
                    Editar
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- vazio -->
        <div v-else class="text-center py-5 text-muted">
          <p class="mb-0">
            {{ users.length === 0 ? 'Nenhum usuário cadastrado.' : 'Nenhum usuário encontrado para essa busca.' }}
          </p>
        </div>

        <!-- footer -->
        <div
          v-if="!loading"
          class="card-footer bg-white border-top py-3 px-4 d-flex justify-content-between align-items-center text-muted small"
        >
          <span>{{ filteredUsers.length }} usuário(s) exibido(s)</span>
          <span v-if="selectedIds.length > 0" class="text-primary fw-medium">
            {{ selectedIds.length }} selecionado(s)
          </span>
        </div>
      </div>
    </div>
  </div>
</template>