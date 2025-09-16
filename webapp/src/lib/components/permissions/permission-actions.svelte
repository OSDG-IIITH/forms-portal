<script lang="ts">
import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
import UserPlusIcon from "@lucide/svelte/icons/user-plus";
import UserMinusIcon from "@lucide/svelte/icons/user-minus";
import ShieldIcon from "@lucide/svelte/icons/shield";
import TrashIcon from "@lucide/svelte/icons/trash";
import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
import * as Command from "$lib/components/ui/command/index.js";
import { Button } from "$lib/components/ui/button/index.js";

interface PermissionActionsProps {
  userId: string;
  userName?: string;
  userEmail?: string;
  roles: Array<"view" | "respond" | "comment" | "analyze" | "edit" | "manage">;
  onPermissionChange?: (userId: string, action: 'add' | 'remove', role: string) => void;
  onRemoveUser?: (userId: string) => void;
}

let { userId, userName, userEmail, roles, onPermissionChange, onRemoveUser }: PermissionActionsProps = $props();

const allRoles = [
  { value: "view", label: "View", description: "Can view the form" },
  { value: "respond", label: "Respond", description: "Can submit responses" },
  { value: "comment", label: "Comment", description: "Can add comments" },
  { value: "analyze", label: "Analyze", description: "Can view responses" },
  { value: "edit", label: "Edit", description: "Can edit the form" },
  { value: "manage", label: "Manage", description: "Full management access" }
];

let open = $state(false);
let addRoleOpen = $state(false);
let removeRoleOpen = $state(false);

function closeAndFocusTrigger() {
  open = false;
  addRoleOpen = false;
  removeRoleOpen = false;
}

function handleAddRole(role: string) {
  onPermissionChange?.(userId, 'add', role);
  closeAndFocusTrigger();
}

function handleRemoveRole(role: string) {
  onPermissionChange?.(userId, 'remove', role);
  closeAndFocusTrigger();
}

function handleRemoveUser() {
  onRemoveUser?.(userId);
  closeAndFocusTrigger();
}

const availableRoles = $derived(
  allRoles.filter(role => !roles.includes(role.value as any))
);
const currentRoles = $derived(
  allRoles.filter(role => roles.includes(role.value as any))
);

function onAddRoleSelect(roleValue: string) {
  handleAddRole(roleValue);
}
function onRemoveRoleSelect(roleValue: string) {
  handleRemoveRole(roleValue);
}
</script>

<DropdownMenu.Root bind:open>
  <DropdownMenu.Trigger>
    {#snippet child({ props })}
      <Button variant="ghost" size="sm" {...props} aria-label="Open menu" class="h-8 w-8 p-0">
        <EllipsisIcon class="h-4 w-4" />
      </Button>
    {/snippet}
  </DropdownMenu.Trigger>
  <DropdownMenu.Content class="w-[220px]" align="end">
    <DropdownMenu.Group>
      <DropdownMenu.Label>Manage Permissions</DropdownMenu.Label>
      
      {#if availableRoles.length > 0}
        <DropdownMenu.Sub bind:open={addRoleOpen}>
          <DropdownMenu.SubTrigger>
            <UserPlusIcon class="mr-2 size-4" />
            Add Permission
          </DropdownMenu.SubTrigger>
          <DropdownMenu.SubContent class="p-0">
            <Command.Root>
              <Command.Input autofocus placeholder="Search permissions..." />
              <Command.List>
                <Command.Empty>No permissions found.</Command.Empty>
                <Command.Group>
                  {#each availableRoles as role}
                    <Command.Item
                      value={role.value}
                      onSelect={() => onAddRoleSelect(role.value)}
                    >
                      <div>
                        <div class="font-medium">{role.label}</div>
                        <div class="text-xs text-muted-foreground">{role.description}</div>
                      </div>
                    </Command.Item>
                  {/each}
                </Command.Group>
              </Command.List>
            </Command.Root>
          </DropdownMenu.SubContent>
        </DropdownMenu.Sub>
      {/if}
      
      {#if currentRoles.length > 0}
        <DropdownMenu.Sub bind:open={removeRoleOpen}>
          <DropdownMenu.SubTrigger>
            <UserMinusIcon class="mr-2 size-4" />
            Remove Permission
          </DropdownMenu.SubTrigger>
          <DropdownMenu.SubContent class="p-0 mx-4">
            <Command.Root>
              <Command.Input autofocus placeholder="Search permissions..." />
              <Command.List>
                <Command.Empty>No permissions found.</Command.Empty>
                <Command.Group>
                  {#each currentRoles as role}
                    <Command.Item
                      value={role.value}
                      onSelect={() => onRemoveRoleSelect(role.value)}
                    >
                      <div class="px-2">
                        <div class="font-medium">{role.label}</div>
                        <div class="text-xs text-muted-foreground">{role.description}</div>
                      </div>
                    </Command.Item>
                  {/each}
                </Command.Group>
              </Command.List>
            </Command.Root>
          </DropdownMenu.SubContent>
        </DropdownMenu.Sub>
      {/if}
      
      <DropdownMenu.Separator />
      <DropdownMenu.Item onSelect={handleRemoveUser}>
        <TrashIcon class="mr-2 size-4" />
        Remove
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>
