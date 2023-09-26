package bastion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type jsonLocalpasswordpolicy struct {
	AllowSameUserAndPassword bool     `json:"allow_same_user_and_password"`
	ID                       string   `json:"id,omitempty"`
	PasswordPolicyName       string   `json:"password_policy_name"`
	PasswordExpiration       int      `json:"password_expiration"`
	PasswordWarningDays      int      `json:"password_warning_days"`
	PasswordMinLength        int      `json:"password_min_length"`
	PasswordMinLowerChars    int      `json:"password_min_lower_chars"`
	PasswordMinUpperChars    int      `json:"password_min_upper_chars"`
	PasswordMinDigitChars    int      `json:"password_min_digit_chars"`
	PasswordMinSpecialChars  int      `json:"password_min_special_chars"`
	LastPasswordsToReject    int      `json:"last_passwords_to_reject"`
	MaxAuthFailures          int      `json:"max_auth_failures"`
	SSHRsaMinLength          int      `json:"ssh_rsa_min_length"`
	ForbiddenPasswords       []string `json:"forbidden_passwords"`
	SSHKeyAlgosAllowed       []string `json:"ssh_key_algos_allowed"`
}

func dataSourceLocalpasswordpolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocalpasswordpolicyRead,
		Schema: map[string]*schema.Schema{
			"password_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"allow_same_user_and_password": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forbidden_passwords": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"last_passwords_to_reject": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_auth_failures": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_expiration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_min_digit_chars": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_min_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_min_lower_chars": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_min_special_chars": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_min_upper_chars": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password_warning_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ssh_key_algos_allowed": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ssh_rsa_min_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLocalpasswordpolicyVersionCheck(version string) error {
	if bchk.InSlice(version, defaultVersionsValid()) {
		return nil
	}

	return fmt.Errorf("data source wallix-bastion_localpasswordpolicy not available with api version %s", version)
}

func dataSourceLocalpasswordpolicyRead(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	c := m.(*Client)
	if err := dataSourceLocalpasswordpolicyVersionCheck(c.bastionAPIVersion); err != nil {
		return diag.FromErr(err)
	}
	cfg, err := readLocalpasswordpolicyOptions(ctx, d.Get("password_policy_name").(string), m)
	if err != nil {
		return diag.FromErr(err)
	}
	fillLocalpasswordpolicy(d, cfg)
	d.SetId(cfg.ID)

	return nil
}

func readLocalpasswordpolicyOptions(
	ctx context.Context, passwordPolicyName string, m interface{},
) (
	jsonLocalpasswordpolicy, error,
) {
	c := m.(*Client)
	body, code, err := c.newRequest(ctx,
		"/localpasswordpolicies/?q=password_policy_name="+passwordPolicyName, http.MethodGet, nil)
	if err != nil {
		return jsonLocalpasswordpolicy{}, err
	}
	if code != http.StatusOK {
		return jsonLocalpasswordpolicy{}, fmt.Errorf("api doesn't return OK: %d with body:\n%s", code, body)
	}
	var results []jsonLocalpasswordpolicy
	err = json.Unmarshal([]byte(body), &results)
	if err != nil {
		return jsonLocalpasswordpolicy{}, fmt.Errorf("unmarshaling json: %w", err)
	}
	if len(results) == 0 {
		return jsonLocalpasswordpolicy{}, fmt.Errorf("password_policy_name %s not found", passwordPolicyName)
	}

	return results[0], nil
}

func fillLocalpasswordpolicy(d *schema.ResourceData, jsonData jsonLocalpasswordpolicy) {
	if tfErr := d.Set("allow_same_user_and_password", jsonData.AllowSameUserAndPassword); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_expiration", jsonData.PasswordExpiration); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_warning_days", jsonData.PasswordWarningDays); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_min_length", jsonData.PasswordMinLength); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_min_lower_chars", jsonData.PasswordMinLowerChars); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_min_upper_chars", jsonData.PasswordMinUpperChars); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_min_digit_chars", jsonData.PasswordMinDigitChars); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("password_min_special_chars", jsonData.PasswordMinSpecialChars); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("last_passwords_to_reject", jsonData.LastPasswordsToReject); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("max_auth_failures", jsonData.MaxAuthFailures); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ssh_rsa_min_length", jsonData.SSHRsaMinLength); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("forbidden_passwords", jsonData.ForbiddenPasswords); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("ssh_key_algos_allowed", jsonData.SSHKeyAlgosAllowed); tfErr != nil {
		panic(tfErr)
	}
}
